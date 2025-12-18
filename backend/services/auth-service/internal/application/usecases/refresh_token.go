package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

var (
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenRequired  = errors.New("refresh token is required")
)

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenUseCase struct {
	userRepo   repositories.UserRepository
	logger     *logger.Logger
	jwtManager *utils.JwtManager
	redis      utils.RedisInterface
}

func NewRefreshTokenUseCase(
	userRepo repositories.UserRepository,
	logger *logger.Logger,
	jwtManager *utils.JwtManager,
	redis utils.RedisInterface,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:   userRepo,
		logger:     logger,
		jwtManager: jwtManager,
		redis:      redis,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}


	claims, err := uc.jwtManager.VerifyToken(input.RefreshToken)
	if err != nil {
		uc.logger.Error("failed to verify refresh token", zap.Error(err))
		return nil, ErrInvalidRefreshToken
	}

	userID, err := uc.redis.GetUserFromRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		uc.logger.Error("refresh token not found in Redis", zap.Error(err))
		return nil, ErrInvalidRefreshToken
	}

	if userID != claims.UserID {
		uc.logger.Error("user ID mismatch between token and Redis", 
			zap.String("token_user_id", claims.UserID),
			zap.String("redis_user_id", userID))
		return nil, ErrInvalidRefreshToken
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user", zap.Error(err))
		return nil, ErrInvalidRefreshToken
	}
	if user == nil {
		uc.logger.Error("user not found", zap.String("user_id", userID))
		return nil, ErrInvalidRefreshToken
	}
	if user.Status != valueobjects.StatusActive {
		uc.logger.Error("user is not active", zap.String("user_id", userID))
		return nil, ErrInvalidRefreshToken
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email.String(),
		user.Role.String(),
		user.Status.String(),
	)
	if err != nil {
		uc.logger.Error("failed to generate access token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	
	newRefreshToken, err := uc.jwtManager.GenerateRefreshToken(
		user.ID,
		user.Email.String(),
		user.Role.String(),
		user.Status.String(),
	)
	if err != nil {
		uc.logger.Error("failed to generate refresh token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}


	if err := uc.redis.StoreAccessToken(ctx, user.ID, accessToken, uc.jwtManager.AccessTokenDuration()); err != nil {
		uc.logger.Warn("failed to store access token", zap.Error(err))
	}

	
	if err := uc.redis.StoreRefreshToken(ctx, user.ID, newRefreshToken, uc.jwtManager.RefreshTokenDuration()); err != nil {
		uc.logger.Warn("failed to store refresh token", zap.Error(err))
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (uc *RefreshTokenUseCase) validateInput(input RefreshTokenInput) error {
	if input.RefreshToken == "" {
		return ErrRefreshTokenRequired
	}
	return nil
}

