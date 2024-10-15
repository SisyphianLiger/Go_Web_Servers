package main

import (
	"encoding/json"
	"net/http"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
)

func (cfg * apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {
    // Process Email Login First
    type parameters struct {
        Email string
        Password string
    }
    p  := parameters{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Inccorrect Body", err)
        return
    }

    // Make Query for sqlc (Find User by Email)
    user, err := cfg.dbc.GetUser(r.Context(), p.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error(), err)
        return
    }
    if invalidLoginErr := auth.CheckPasswordHash(user.HashedPassword, p.Password); invalidLoginErr != nil {
        respondWithError(w, http.StatusUnauthorized, "Not a valid password, please try again", err)
        return
    }
    
    // Respond
    respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
    })
}
