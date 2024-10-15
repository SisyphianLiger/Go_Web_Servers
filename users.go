package main

import (
	"encoding/json"
	"net/http"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
)

// UPGRADE
func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	p := params{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Incorrect Json Key at Body", err)
		return
	}
	// Here we need to add the password	
	hashedPassword, err := auth.HashPassword(p.Password)	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password Hash did not work please check length", err)
		return
	}
	user, uerr := cfg.dbc.CreateUser(r.Context(), database.CreateUserParams{Email: p.Email, HashedPassword: hashedPassword})
	if uerr != nil {
		respondWithError(w, http.StatusInternalServerError,
			"User with that email already exits", err)
		return
	}

	respondWithJson(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})

}
