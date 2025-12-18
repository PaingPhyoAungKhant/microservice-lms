package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
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
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailRequired = errors.New("email is required")
	ErrUsernameRequired = errors.New("username is required")
	ErrPasswordRequired = errors.New("password is required")
) 

type RegisterStudentInput struct {
	Email string 
	Username string
	Password string
}

type RegisterStudentUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	rabbitMQConfig *config.RabbitMQConfig
	logger *logger.Logger
	redis utils.RedisInterface
	apiGatewayURL string
}

func NewRegisterStudentUseCase(
	userRepo repositories.UserRepository,
	 publisher messaging.Publisher,
	 logger *logger.Logger,
	 rabbitMQConfig *config.RabbitMQConfig,
	 redis utils.RedisInterface,
	 apiGatewayURL string,
	) *RegisterStudentUseCase {
		return &RegisterStudentUseCase{
			userRepo: userRepo,
			publisher: publisher,
			logger: logger,
			rabbitMQConfig: rabbitMQConfig,
			redis: redis,
			apiGatewayURL: apiGatewayURL,
		}
}

func (uc *RegisterStudentUseCase)Execute(
	ctx context.Context, 
	input RegisterStudentInput,
	) (*dtos.UserDTO, error) {
	validationErr := uc.validateInput(input)
	if validationErr != nil {
		return nil, validationErr
	}

	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
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

	user := entities.NewUser(email, input.Username, valueobjects.RoleStudent, passwordHash)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	token := uuid.New().String()
	err = uc.redis.StoreVerifyEmailToken(ctx, user.ID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to store verify email token: %w", err)
	}

	emailVerificationURL := utils.GenerateEmailVerificationURL(uc.apiGatewayURL, token)

	event := events.AuthStudentRegisteredEvent{
		ID:                  user.ID,
		Email:               user.Email.String(),
		Username:            user.Username,
		Role:                user.Role.String(),
		Status:              user.Status.String(),
		EmailVerified:       user.EmailVerified,
		EmailVerifiedAt:     user.EmailVerifiedAt,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		EmailVerificationURL: emailVerificationURL,
	}

	if err := uc.publisher.Publish(ctx, events.EventTypeAuthStudentRegistered, event); err != nil {
		uc.logger.Error("failed to publish student registered event", zap.Error(err))
	}

	var dto dtos.UserDTO
	dto.FromEntity(user)
	return &dto, nil
}

func (uc *RegisterStudentUseCase) validateInput(input RegisterStudentInput) error {
	if input.Email == "" {
		return  ErrEmailRequired
	}

	if input.Username == "" {
		return  ErrUsernameRequired
	}
	if input.Password == "" {
		return  ErrPasswordRequired
	}

	return nil
}