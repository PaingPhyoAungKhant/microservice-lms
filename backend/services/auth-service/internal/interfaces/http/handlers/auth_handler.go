package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	usecases "github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
)

type AuthHandler struct {
	loginUseCase *usecases.LoginUseCase
	registerStudentUseCase *usecases.RegisterStudentUseCase
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase
	resetPasswordUseCase *usecases.ResetPasswordUseCase
	verifyOTPUseCase *usecases.VerifyOTPUseCase
	verifyUseCase *usecases.VerifyUseCase
	verifyEmailUseCase *usecases.VerifyEmailUseCase
	requestEmailVerifyUseCase *usecases.RequestEmailVerifyUseCase
	refreshTokenUseCase *usecases.RefreshTokenUseCase
}

func NewAuthHandler(
	loginUseCase *usecases.LoginUseCase,
	registerStudentUseCase *usecases.RegisterStudentUseCase,
	forgotPasswordUseCase *usecases.ForgotPasswordUseCase,
	resetPasswordUseCase *usecases.ResetPasswordUseCase,
	verifyOTPUseCase *usecases.VerifyOTPUseCase,
	verifyUseCase *usecases.VerifyUseCase,
	verifyEmailUseCase *usecases.VerifyEmailUseCase,
	requestEmailVerifyUseCase *usecases.RequestEmailVerifyUseCase,
	refreshTokenUseCase *usecases.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
		registerStudentUseCase: registerStudentUseCase,
		forgotPasswordUseCase: forgotPasswordUseCase,
		resetPasswordUseCase: resetPasswordUseCase,
		verifyOTPUseCase: verifyOTPUseCase,
		verifyUseCase: verifyUseCase,
		verifyEmailUseCase: verifyEmailUseCase,
		requestEmailVerifyUseCase: requestEmailVerifyUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"Password@123"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.LoginInput{
		Email:     req.Email,
		Password:   req.Password,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	output, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send OTP to user's email for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} map[string]interface{} "OTP sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or email not found"
// @Router /forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.ForgotPasswordInput{
		Email: req.Email,
	}

	output, err := h.forgotPasswordUseCase.Execute(c.Request.Context(), input)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
	OTP   string `json:"otp" binding:"required" example:"123456"`
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify the OTP sent to user's email for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Email and OTP"
// @Success 200 {object} map[string]interface{} "OTP verified successfully, returns reset token"
// @Failure 400 {object} map[string]interface{} "Invalid OTP or request body"
// @Router /verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.VerifyOTPInput{
		Email: req.Email,
		OTP: req.OTP,
	}


	output, err := h.verifyOTPUseCase.Execute(c.Request.Context(), input)
	if err != nil || !output.IsValid {
		c.JSON(400, gin.H{"error": output.ErrorMessage, "details": err.Error()})
		return
	}

	c.JSON(200, output)
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"reset-token-here"`
	NewPassword string `json:"new_password" binding:"required" example:"Password@123"`
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using the token received from OTP verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid token or request body"
// @Router /reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.ResetPasswordInput{
		Token:    req.Token,
		NewPassword: req.NewPassword,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	output, err := h.resetPasswordUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

// Verify godoc
// @Summary Verify JWT token
// @Description Verify JWT access token and optionally check for required role. Returns user information in headers.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param X-Required-Role header string false "Required role for access"
// @Success 200 "Token is valid"
// @Failure 400 {object} map[string]interface{} "Token is required"
// @Failure 401 "Token is invalid or expired"
// @Failure 403 "Insufficient permissions"
// @Router /verify [get]
func (h *AuthHandler) Verify(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(400, gin.H{"error": "Token is required"})
		return
	}

	token := authHeader
	authHeaderLower := strings.ToLower(authHeader)
	if strings.HasPrefix(authHeaderLower, "bearer ") {
		token = strings.TrimSpace(authHeader[7:])
	}

	if token == "" {
		c.JSON(400, gin.H{"error": "Token is required"})
		return
	}

	input := usecases.VerifyInput{
		Token: token,
		RequiredRole: c.GetHeader("X-Required-Role"),
	}

	user, err := h.verifyUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch err {
		case usecases.ErrInsufficientPermissions:
		  c.Status(403)
			return
		default:
			c.Status(401)
			return
		}
	}

	c.Header("X-User-ID", user.ID)
	c.Header("X-User-Email", user.Email)
	c.Header("X-User-Role", user.Role)
	c.Status(200)
}

type RegisterStudentRequest struct {
	Email    string `json:"email" binding:"required" example:"student@example.com"`
	Username string `json:"username" binding:"required" example:"student123"`
	Password string `json:"password" binding:"required" example:"Password@123"`
}

// RegisterStudent godoc
// @Summary Register new student
// @Description Register a new student account. An email verification link will be sent to the provided email.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterStudentRequest true "Student registration details"
// @Success 201 {object} map[string]interface{} "Student registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or email/username already exists"
// @Router /register [post]
func (h *AuthHandler) RegisterStudent(c *gin.Context) {
	var req RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.RegisterStudentInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	output, err := h.registerStudentUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, output)
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user's email address using the token from verification link
// @Tags auth
// @Produce json
// @Param token query string true "Email verification token"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid token"
// @Router /verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "token is required"})
		return
	}

	input := usecases.VerifyEmailInput{
		Token: token,
	}

	output, err := h.verifyEmailUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

type RequestEmailVerifyRequest struct {
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

// RequestEmailVerify godoc
// @Summary Request email verification
// @Description Request a new email verification link to be sent to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RequestEmailVerifyRequest true "Email address"
// @Success 200 {object} map[string]interface{} "Verification email sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or email not found"
// @Router /request-email-verify [post]
func (h *AuthHandler) RequestEmailVerify(c *gin.Context) {
	var req RequestEmailVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.RequestEmailVerifyInput{
		Email: req.Email,
	}

	output, err := h.requestEmailVerifyUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "New access token generated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid or expired refresh token"
// @Router /refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	input := usecases.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}

	output, err := h.refreshTokenUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}

// Health godoc
// @Summary Health check
// @Description Check if the auth service is running and healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func (h *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
	})
}