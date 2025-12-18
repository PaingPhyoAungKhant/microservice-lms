// Package usecases
package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailRequired = errors.New("email is required")
	ErrUsernameRequired = errors.New("username is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrRoleRequired = errors.New("role is required")
	)

type CreateUserInput struct {
	Email    string
	Username string
	Password string
	Role     string
}

type CreateUserUseCase struct {
	userRepo  repositories.UserRepository
	publisher messaging.Publisher
	logger    *logger.Logger
	redis utils.RedisInterface
	apiGatewayURL string
}

func NewCreateUserUseCase(userRepo repositories.UserRepository, publisher messaging.Publisher, logger *logger.Logger, redis utils.RedisInterface, apiGatewayURL string) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:  userRepo,
		publisher: publisher,
		logger:    logger,
		redis: redis,
		apiGatewayURL: apiGatewayURL,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*dtos.UserDTO, error) {
	
	validationErr := uc.validateInput(input)
	if validationErr != nil {
		return nil, validationErr
	}

	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	role, err := valueobjects.NewRole(input.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	existingEmail, _ := uc.userRepo.FindByEmail(ctx, email.String())
	if existingEmail != nil {
		return nil, ErrEmailAlreadyExists
	}

	existingUsername, _ := uc.userRepo.FindByUsername(ctx, input.Username)
	if existingUsername != nil {
		return nil, ErrUsernameAlreadyExists
	}

	if err := utils.ValidatePassword(input.Password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := entities.NewUser(email, input.Username, role, passwordHash)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token := uuid.New().String()
	err = uc.redis.StoreVerifyEmailToken(ctx, user.ID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verify email token: %w", err)
	}
	
	emailVerificationURL := utils.GenerateEmailVerificationURL(uc.apiGatewayURL, token)
	event := events.UserCreatedEvent{
		ID:         user.ID,
		Email:      user.Email.String(),
		Username:   user.Username,
		Role:       user.Role.String(),
		Status:     user.Status.String(),
		EmailVerified: user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		EmailVerificationURL: emailVerificationURL,
	}

	if uc.publisher != nil {
		if err := uc.publisher.Publish(ctx, events.EventTypeUserCreated, event); err != nil {
			uc.logger.Error("Failed to publish user created event", zap.Error(err))
		}
	}

	uc.logger.Info("User Created Successfully.",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email.String()),
	)

	var dto dtos.UserDTO
	dto.FromEntity(user)
	return &dto, nil
}

func (uc *CreateUserUseCase) validateInput(input CreateUserInput) error {
	if input.Email == "" {
		return  ErrEmailRequired
	}

	if input.Username == "" {
		return  ErrUsernameRequired
	}

	if input.Role == "" {
		return  ErrRoleRequired
	}

	if input.Password == "" {
		return  ErrPasswordRequired
	}	

	return nil
}