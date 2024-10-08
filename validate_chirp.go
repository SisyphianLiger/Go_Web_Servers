package main

import (
	"encoding/json"
	"net/http"
)


func validateChirp(w http.ResponseWriter, r *http.Request){
    // Chirp that takes in Body
    type Chirp struct {
        Chirp string `json:"body"`
    }

    type validResponse struct {
        Response bool `json:"valid"`
    }

    chirp := Chirp{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&chirp)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't code parameters", err)
        return 
    }

    if len(chirp.Chirp) > 140 {
        respondWithError(w, http.StatusBadRequest,"Chirp is too long", nil) 
        return 
    }

    respondWithJson(w, http.StatusOK, validResponse{
        Response: true, 
    })    
    
}

