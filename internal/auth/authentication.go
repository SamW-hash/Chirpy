package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashStr, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Couldn't hash password: %w", err)
	}
	return hashStr, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Couldn't compare password: %w", err)
	}
	return result, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing token: %w", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	registeredClaims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &registeredClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error validating: %w", err)
	}
	id, err := uuid.Parse(registeredClaims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing id: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if len(bearerToken) == 0 {
		return "", fmt.Errorf("No bearer token")
	}
	return strings.TrimPrefix(bearerToken, "Bearer "), nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	rToken := hex.EncodeToString(token)
	return rToken
}

func GetAPIKey(headers http.Header) (string, error) {
	APIKey := headers.Get("Authorization")
	if len(APIKey) == 0 {
		return "", fmt.Errorf("No API Key")
	}
	return strings.TrimPrefix(APIKey, "ApiKey "), nil
}
