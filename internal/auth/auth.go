package auth
import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err:= bcrypt.GenerateFromPassword([]byte(password), len(password))
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))	
	if err != nil {
		return err
	}
	return nil
}
