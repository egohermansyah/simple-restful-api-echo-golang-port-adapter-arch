package generatepassword

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string, passwordSalt string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s-%s", password, passwordSalt)), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
