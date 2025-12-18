package utils

import (
	"fmt"
	"math"
	"math/rand"
)

func GenerateOTP(length int) string {
	if length <= 0 {
		length = 6
	}

	max := int(math.Pow10(length))

	min := int(math.Pow10(length - 1))

	otp := rand.Intn(max-min) + min

	return fmt.Sprintf("%0*d", length, otp)
}