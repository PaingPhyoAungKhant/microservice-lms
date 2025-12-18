package templates

import _ "embed"

//go:embed email_verification.html
var EmailVerification string

//go:embed forgot_password_otp.html
var ForgotPasswordOTP string

