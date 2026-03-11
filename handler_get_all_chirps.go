package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")
	if sortDirection == "desc" {
		sortDirection = "desc"
	} else {
		sortDirection = "asc"
	}
	if len(s) != 0 {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author id", err)
			return
		}
		chirps, err := cfg.db.GetChirpsByUser(r.Context(), authorID)
		type chirpResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}
		sort.Slice(chirps, func(i, j int) bool {
			if sortDirection == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
		respBody := make([]chirpResponse, len(chirps))
		for i, c := range chirps {
			respBody[i] = chirpResponse{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			}
		}
		respondWithJSON(w, 200, respBody)
		return
	}
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
	}
	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
	type chirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	respBody := make([]chirpResponse, len(chirps))
	for i, c := range chirps {
		respBody[i] = chirpResponse{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	respondWithJSON(w, 200, respBody)
}
