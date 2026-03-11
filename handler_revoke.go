package main

import (
	"net/http"
	"workspace/sam/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing validation token", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}
	respondWithJSON(w, 204, nil)
}
