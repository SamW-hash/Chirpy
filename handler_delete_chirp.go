package main

import (
	"net/http"
	"workspace/sam/Chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	//Who is the person making this request?
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid token provided", err)
		return
	}
	user, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Token not valid", err)
		return
	}

	//What are they requesting to be deleted?
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirp id", err)
		return
	}

	//Did they make the chirp they want to delete?
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if chirp.UserID != user {
		respondWithError(w, 403, "Unauthorized action", err)
		return
	}

	//Delete it
	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	respondWithJSON(w, 204, nil)

}
