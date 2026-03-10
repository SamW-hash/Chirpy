package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "my-secret"
	userID := uuid.New()
	expiresIn := time.Hour

	// 1. Create the token
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("couldn't make jwt: %v", err)
	}

	// 2. Validate the token
	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("couldn't validate jwt: %v", err)
	}

	// 3. Check if the ID matches
	if returnedID != userID {
		t.Errorf("expected %v, got %v", userID, returnedID)
	}
}
