package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset not allowed on this platform", nil)
		return
	}
	if err := cfg.db.ResetDatabase(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset database", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database Reset!"))
}
