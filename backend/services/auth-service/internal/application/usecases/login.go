package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)


var (
	ErrUserNotFound = errors.New("user not found")
)


type LoginInput struct {
	Email     string
	Password  string
	IPAddress string 
	UserAgent string
}

type LoginOutput struct {
	User         *dtos.UserDTO `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
}

type LoginUseCase struct {
	userRepo repositories.UserRepository
	publisher messaging.Publisher
	logger *logger.Logger
	jwtManager *utils.JwtManager
	redis utils.RedisInterface
}

func NewLoginUseCase(
	userRepo repositories.UserRepository,
	publisher messaging.Publisher,
	logger *logger.Logger,
	jwtManager *utils.JwtManager,
	redis utils.RedisInterface,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		publisher: publisher,
		logger: logger,
		jwtManager: jwtManager,
		redis: redis,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	validationErr := uc.validateInput(input)
	if validationErr != nil {
		return nil, validationErr
	}

	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := utils.VerifyPassword(input.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidPassword
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID, user.Email.String(), user.Role.String(), user.Status.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	sessionID := uuid.NewString()


	if err := uc.redis.StoreAccessToken(ctx, user.ID, accessToken, uc.jwtManager.AccessTokenDuration()); err != nil {
		uc.logger.Warn("failed to store access token", zap.Error(err))
	}

	if err := uc.redis.StoreRefreshToken(ctx, user.ID, refreshToken, uc.jwtManager.RefreshTokenDuration()); err != nil {
		uc.logger.Warn("failed to store refresh token", zap.Error(err))
	}

	
	if err := uc.redis.StoreUserSession(
		ctx,
		sessionID,
		user.ID,
		accessToken,
		refreshToken,
		input.IPAddress,
		input.UserAgent,
		uc.jwtManager.AccessTokenDuration(),
	); err != nil {
		uc.logger.Warn("failed to store user session", zap.Error(err))
	}

	event := events.AuthUserLoggedInEvent{
		ID:            user.ID,
		Email:         user.Email.String(),
		Username:      user.Username,
		Role:          user.Role.String(),
		Status:        user.Status.String(),
		EmailVerified: user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	
	if err := uc.publisher.Publish(ctx, events.EventTypeAuthUserLoggedIn, event); err != nil {
		uc.logger.Error("failed to publish user logged in event", zap.Error(err))
	}

	var dto dtos.UserDTO
	dto.FromEntity(user)

	return &LoginOutput{
		User:         &dto,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *LoginUseCase) validateInput(input LoginInput) error {
	if input.Email == "" {
		return  ErrEmailRequired
	}


	if input.Password == "" {
		return  ErrPasswordRequired
	}

	return nil
}