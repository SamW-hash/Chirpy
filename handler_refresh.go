package main

import (
	"net/http"
	"time"
	"workspace/sam/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing validation token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not valid", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}
	type returnVals struct {
		Token string `json:"token"`
	}
	respBody := returnVals{
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
