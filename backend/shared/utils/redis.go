package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
	"github.com/redis/go-redis/v9"
)


const (
	RedisKeyAccessToken        = "auth:token:access:%s"
	RedisKeyRefreshToken       = "auth:token:refresh:%s"
	RedisKeySession            = "auth:session:%s"
	RedisKeyForgotPasswordOTP = "auth:forgot_password_otp:%s"
	RedisKeyResetPassword      = "auth:reset_password:%s"
	RedisKeyVerifyEmail        = "auth:verify_email:%s"
)

type RedisInterface interface {
	StoreAccessToken(ctx context.Context, userID string, accessToken string, expiration time.Duration) error
	GetUserFromAccessToken(ctx context.Context, accessToken string) (string, error)
	RevokeAccessToken(ctx context.Context, accessToken string) error
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiration time.Duration) error
	GetUserFromRefreshToken(ctx context.Context, refreshToken string) (string, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	StoreUserSession(ctx context.Context, sessionID, userID, accessToken, refreshToken, ipAddress, userAgent string, expiration time.Duration) error
	GetUserSession(ctx context.Context, sessionID string) (*SessionData, error)
	DeleteUserSession(ctx context.Context, sessionID string) error
	ExpireUserSession(ctx context.Context, sessionID string, expiration time.Duration) error
	UpdateLastActivity(ctx context.Context, sessionID string) error
	StoreForgotPasswordOTP(ctx context.Context, userID, otp string) error
	GetUserFromForgotPasswordOTP(ctx context.Context, otp string) (string, error)
	RevokeForgotPasswordOTP(ctx context.Context, otp string) error
	StoreResetPasswordToken(ctx context.Context, userID, token string) error
	GetUserFromResetPasswordToken(ctx context.Context, token string) (string, error)
	RevokeResetPasswordToken(ctx context.Context, token string) error
	StoreVerifyEmailToken(ctx context.Context, userID, token string) error
	GetUserFromVerifyEmailToken(ctx context.Context, token string) (string, error)
	RevokeVerifyEmailToken(ctx context.Context, token string) error
}

type Redis struct {
	client *redis.Client
}

func NewRedis(config *config.RedisConfig) (*Redis, error) {
  addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Println("Redis host:", config.Host)
	log.Println("Redis port:", config.Port)
	log.Println("Redis address:", addr)
	log.Println("Redis password:", config.Password)
	log.Println("Redis DB:", config.DB)
	client := redis.NewClient(
		&redis.Options{
			Addr: addr,
			Password: config.Password,
			DB: config.DB,
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)	
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &Redis{client}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Set(
	ctx context.Context, 
	key string, 
	value interface {}, 
	expiration time.Duration,
) error {
	var data []byte 
	var err error 

	switch v := value.(type) {
	case string: 
	  data = []byte(v)
	case []byte:
		data = v 
	default: 
	  data, err = json.Marshal(v)
	  if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	  }
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (r *Redis) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if key %s exists: %w", key, err)
	}
	return exists > 0, nil
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to expire key %s: %w", key, err)
	}
	return nil
}

func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	result, err := r.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return result, nil
}

/*	------------------------------------------- Auth Token Management ------------------------------------------- */

func (r *Redis) StoreAccessToken(ctx context.Context, userID string, accessToken string, expiration time.Duration) error {
 key := fmt.Sprintf(RedisKeyAccessToken, accessToken)
 return r.Set(ctx, key, userID, expiration) 
}

func (r *Redis) GetUserFromAccessToken(ctx context.Context, accessToken string ) (string, error) {
	key := fmt.Sprintf(RedisKeyAccessToken, accessToken)
	userID, err := r.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get user from access token: %w", err)
	}
	return userID, nil
}

func (r *Redis) RevokeAccessToken(ctx context.Context, accessToken string) error {
	key := fmt.Sprintf(RedisKeyAccessToken, accessToken)
	return r.Delete(ctx, key)
}

func (r *Redis) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expiration time.Duration) error {
	key := fmt.Sprintf(RedisKeyRefreshToken, refreshToken)
	return r.Set(ctx,key, userID, expiration)
}

func (r *Redis) GetUserFromRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := fmt.Sprintf(RedisKeyRefreshToken, refreshToken)
	userID, err := r.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get user from refresh token: %w", err)
	}
	return userID, nil
}

func (r *Redis) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(RedisKeyRefreshToken, refreshToken)
	return r.Delete(ctx, key)
}

/*	------------------------------------------- User Session Management ------------------------------------------- */

type SessionData struct {
	UserID       string    `json:"user_id"`
	SessionID    string    `json:"session_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    int64     `json:"created_at"`
	ExpiresAt    int64     `json:"expires_at"`
	LastActivity int64     `json:"last_activity"`
}


func (r *Redis) StoreUserSession(ctx context.Context, sessionID, userID, accessToken, refreshToken, ipAddress, userAgent string, expiration time.Duration) error {
	key := fmt.Sprintf(RedisKeySession, sessionID)
	now := time.Now()
	sessionData := SessionData{
		UserID:       userID,
		SessionID:    sessionID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    now.Unix(),
		ExpiresAt:    now.Add(expiration).Unix(),
		LastActivity: now.Unix(),
	}
	return r.Set(ctx, key, sessionData, expiration)
}


func (r *Redis) GetUserSession(ctx context.Context, sessionID string) (*SessionData, error) {
	key := fmt.Sprintf(RedisKeySession, sessionID)
	var sessionData SessionData
	if err := r.GetJSON(ctx, key, &sessionData); err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}
	return &sessionData, nil
}

func (r *Redis) DeleteUserSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf(RedisKeySession, sessionID)
	return r.Delete(ctx, key)
}

func (r *Redis) ExpireUserSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	key := fmt.Sprintf(RedisKeySession, sessionID)
	return r.Expire(ctx, key, expiration)
}

func (r *Redis) UpdateLastActivity(ctx context.Context, sessionID string) error {
	session, err := r.GetUserSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for activity update: %w", err)
	}
	
	session.LastActivity = time.Now().Unix()
	key := fmt.Sprintf(RedisKeySession, sessionID)
	return r.Set(ctx, key, session, time.Until(time.Unix(session.ExpiresAt, 0)))
}

/*	------------------------------------------- Forgot Password OTP Management ------------------------------------------- */

func (r *Redis) StoreForgotPasswordOTP(ctx context.Context, userID, otp string) error {
	key := fmt.Sprintf(RedisKeyForgotPasswordOTP, otp)
	return r.Set(ctx, key, userID, 15*time.Minute)
}

func (r *Redis) GetUserFromForgotPasswordOTP(ctx context.Context, otp string) (string, error) {
	key := fmt.Sprintf(RedisKeyForgotPasswordOTP, otp)
	return r.Get(ctx, key)
}

func (r *Redis) RevokeForgotPasswordOTP(ctx context.Context, otp string) error {
	key := fmt.Sprintf(RedisKeyForgotPasswordOTP, otp)
	return r.Delete(ctx, key)
}
/*	------------------------------------------- Reset Password Management ------------------------------------------- */

func (r *Redis) StoreResetPasswordToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf(RedisKeyResetPassword, token)
	return r.Set(ctx, key, userID, 15*time.Minute)
}

func (r *Redis) GetUserFromResetPasswordToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf(RedisKeyResetPassword, token)
	return r.Get(ctx, key)
}

func (r *Redis) RevokeResetPasswordToken(ctx context.Context, token string) error {
	key := fmt.Sprintf(RedisKeyResetPassword, token)
	return r.Delete(ctx, key)
}
/*	------------------------------------------- Verify Email Management ------------------------------------------- */

func (r *Redis) StoreVerifyEmailToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf(RedisKeyVerifyEmail, token)
	return r.Set(ctx, key, userID, 24*time.Hour)
}

func (r *Redis) GetUserFromVerifyEmailToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf(RedisKeyVerifyEmail, token)
	return r.Get(ctx, key)
}

func (r *Redis) RevokeVerifyEmailToken(ctx context.Context, token string) error {
	key := fmt.Sprintf(RedisKeyVerifyEmail, token)
	return r.Delete(ctx, key)
}