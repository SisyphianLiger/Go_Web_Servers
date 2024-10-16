package auth
import (
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",	
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	retToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Token Secret may be invalid please check again")
	}

	return retToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type claimParameters struct {
		jwt.RegisteredClaims
	}
	// Something Here
	token,err := jwt.ParseWithClaims(tokenString, &claimParameters{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	// Returns nil if token errors
	if err != nil {
		return uuid.Nil, err
	}

	// Checks if in GetSubject	
	claimID, cError := token.Claims.GetSubject()
	if cError != nil {
		return uuid.Nil, cError
	}

	// Checks if Parsable
	result, err := uuid.Parse(claimID)
	if err != nil {
		return uuid.Nil, err
	}
	// Ignore return value just do ide doesn't error
	return result, nil
}
