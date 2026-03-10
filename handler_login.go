package main

import (
	"encoding/json"
	"net/http"
	"time"
	"workspace/sam/Chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid email/password", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	validation, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !validation {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	var expirationSeconds int
	if params.ExpiresInSeconds == nil {
		expirationSeconds = 3600
	} else if *params.ExpiresInSeconds > 3600 {
		expirationSeconds = 3600
	} else {
		expirationSeconds = *params.ExpiresInSeconds
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expirationSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	respBody := returnVals{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, respBody)

}
