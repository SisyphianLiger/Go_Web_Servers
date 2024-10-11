package main

import (
	"net/http"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) resetUserTable(w http.ResponseWriter, r *http.Request) {

	if cfg.dev != "dev" {
		respondWithError(w, http.StatusForbidden,
			"Environment Variable not set to dev", nil)
		return
	}


	if err := cfg.dbc.ResetAllUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Table not Found: Check if users table exists", err)
		return
	}
	

	respondWithJson(w, http.StatusOK, nil)
}

