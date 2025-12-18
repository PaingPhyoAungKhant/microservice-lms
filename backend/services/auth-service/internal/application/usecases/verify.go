package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/paingphyoaungkhant/asto-microservice/shared/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)

type VerifyUseCase struct {
	userRepo   repositories.UserRepository
	publisher  messaging.Publisher
	logger     *logger.Logger
	jwtManager *utils.JwtManager
	redis      utils.RedisInterface
}

type VerifyInput struct {
	Token        string
	RequiredRole string
}



func NewVerifyUseCase(userRepo repositories.UserRepository, publisher messaging.Publisher, logger *logger.Logger, jwtManager *utils.JwtManager, redis utils.RedisInterface) *VerifyUseCase {
	return &VerifyUseCase{
		userRepo:   userRepo,
		publisher:  publisher,
		logger:     logger,
		jwtManager: jwtManager,
		redis:      redis,
	}
}

func (uc *VerifyUseCase) Execute(ctx context.Context, input VerifyInput) (*dtos.UserDTO, error) {
	claims, err := uc.jwtManager.VerifyToken(input.Token)
	if err != nil {
		uc.logger.Error("failed to verify token", zap.Error(err))
		return nil, ErrUnauthorized
	}

	userID, err := uc.redis.GetUserFromAccessToken(ctx, input.Token)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to find user", zap.Error(err))
		return nil, ErrUnauthorized
	}
	if user == nil {
		uc.logger.Error("user not found", zap.String("user_id", claims.UserID))
		return nil, ErrUnauthorized
	}
	if user.Status != valueobjects.StatusActive {
		uc.logger.Error("user is not active", zap.String("user_id", claims.UserID))
		return nil, ErrUnauthorized
	}

	if input.RequiredRole != "" {
		if !uc.hasRequiredRole(claims.Role, input.RequiredRole) {
			return nil, ErrInsufficientPermissions
		}
		if !uc.hasRequiredRole(user.Role.String(), input.RequiredRole) {
		return nil, ErrInsufficientPermissions
		}
	}

	dto := dtos.UserDTO{}
	dto.FromEntity(user)
	return &dto, nil
}

func (uc *VerifyUseCase) hasRequiredRole(userRole, requiredRoles string) bool {
	if requiredRoles == "" {
		return true
	}

	roles := strings.Split(requiredRoles, ",")
	for _, role := range roles {
		if strings.TrimSpace(role) == userRole {
			return true
		}
	}
	return false
}

