package main

import (
	"encoding/json"
	"net/http"
	"time"
	"workspace/sam/Chirpy/internal/auth"
	"workspace/sam/Chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token provided", err)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not valid", err)
		return
	}
	hashPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	updUser, err := cfg.db.EditUser(r.Context(), database.EditUserParams{
		Email:          params.Email,
		HashedPassword: hashPass,
		ID:             id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	respBody := returnVals{
		ID:          id,
		CreatedAt:   updUser.CreatedAt,
		UpdatedAt:   updUser.UpdatedAt,
		Email:       updUser.Email,
		IsChirpyRed: updUser.IsChirpyRed.Bool,
	}
	respondWithJSON(w, http.StatusOK, respBody)

}
