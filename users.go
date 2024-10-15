package main

import (
	"encoding/json"
	"net/http"
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
	
	user, err := cfg.dbc.CreateUser(r.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"User with that email already exits", err)
		return
	}


	respondWithJson(w, http.StatusCreated, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
