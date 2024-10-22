package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",	
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
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


func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization Header Not Found")
	}
	
	authSlice := strings.SplitN(authHeader, " ", 2)
	if strings.ToLower(authSlice[0]) != "bearer" {
		return "", fmt.Errorf("Header did not have Bearer Token but %s", authSlice[0])
	}
	return strings.TrimSpace(authSlice[1]), nil
}


func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, hexError := rand.Read(b)
	if hexError != nil {
		return "", fmt.Errorf("Problem allocating 256-bit byte array to Random: %v", hexError)
	}
	result := hex.EncodeToString(b)
	return result, nil

}


func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization Header Not Found")
	}
	
	authSlice := strings.SplitN(authHeader, " ", 2)
	if len(authSlice) != 2 || strings.ToLower(authSlice[0]) != "apikey" {
			return "", fmt.Errorf("Invalid Authorization Header format")
}
	return strings.TrimSpace(authSlice[1]), nil

}
