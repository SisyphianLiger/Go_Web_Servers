package auth

import (
	"testing"
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
