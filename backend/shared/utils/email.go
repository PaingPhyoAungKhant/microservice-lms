package utils

import "fmt"

func GenerateEmailVerificationURL(apiGatewayURL, token string) string {
	return fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", apiGatewayURL, token)
}

