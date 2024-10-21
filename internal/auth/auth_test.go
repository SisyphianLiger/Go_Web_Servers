package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func CheckAndPassFunction(t *testing.T) {
	password := "Hello1235%^*#JAL><"
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Password is not able to be hashed: %s, Error: %v", password, err)
	}
	if err := CheckPasswordHash(password, hashedPass); err != nil {
		t.Fatalf("Password and Hashed Password do not Match: %v", err)	
	}
}

func CheckThatIncorrectPasswordFails(t * testing.T) {
	password := "Hello1235%^*#JAL><"
	badPassword := "Hello1235%^*#JAL<!"
	hp1, _ := HashPassword(password)
	hp2, _ := HashPassword(badPassword)
	if hp1 == hp2 {
		t.Fatalf("Hashing Has failed, Password: %s and BadPassword: %s are the same", password, badPassword)
	}
}

func TestJWTCreationAndValidation(t *testing.T) {
	
	// Test Data
	userID := uuid.New()
	secret := "test_secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Unable to Make JWT")
	}

	vUser, vErr := ValidateJWT(token, secret)
	if vErr != nil {
		t.Fatalf("Validation Created error: %v", vErr)
	}

	if userID != vUser {
		t.Errorf("Uuid for new user and Validated ID do not match")
	}

}

func TestJWTDetectsExpiredToken(t *testing.T) {
	
	// Test Data
	userID := uuid.New()
	secret := "test_secret"

	token, err := MakeJWT(userID, secret, time.Nanosecond)
	if err != nil {
		t.Fatalf("Unable to Make JWT")
	}
	time.Sleep(time.Millisecond)

	_, vErr := ValidateJWT(token, secret)
	if vErr == nil {
		t.Fatalf("Problem: Token should be expired but is in use")
	}

}

func TestJWTWrongSecret(t *testing.T) {
	
	// Test Data
	userID := uuid.New()
	secret := "test_secret"
	wrongSecret := "borkingthecode"

	token, err := MakeJWT(userID, secret, time.Nanosecond)
	if err != nil {
		t.Fatalf("Unable to Make JWT")
	}
	time.Sleep(time.Millisecond)

	_, vErr := ValidateJWT(token, wrongSecret)
	if vErr == nil {
		t.Fatalf("Problem: Use of Wrong Secret did not cause errors")
	}
}

func TestGetBearerTokenSuccess(t *testing.T) {
	header := make(http.Header)	
	header.Set("Authorization", "Bearer LIGHT_WEIGHT")
	token, error := GetBearerToken(header)
	if error != nil {
		t.Errorf("Got %s instead of LIGHT_WEIGHT", token)
	}
}

func TestNoAuthorizationHeader(t *testing.T) {
	header := make(http.Header)	
	header.Set("", "Bearer LIGHT_WEIGHT")
	_, error := GetBearerToken(header)
	if error == nil {
		t.Fatalf("No Header Given but still successful, please check function specifcation")
	}
}
