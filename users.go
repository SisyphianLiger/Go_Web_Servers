package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
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

	log.Printf("Request body right before sending: %+v", user.ID)

	respondWithJson(w, http.StatusCreated, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
