package utils

import (
	"math/rand"
)

func GenOtp() string {
	characters := "abcdefghijklmnopqrstuvwxyz0123456789"

	otpLength := 4
	otp := make([]byte, otpLength)

	for i := range otpLength {
		index := rand.Intn(len(characters))
		otp[i] = characters[index]
	}

	return string(otp)
}
