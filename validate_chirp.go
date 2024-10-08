package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func validateChirp(w http.ResponseWriter, r *http.Request){
    // Chirp that takes in Body
    type Chirp struct {
        Chirp string `json:"body"`
    }

    type validResponse struct {
        Response string `json:"cleaned_body"`
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
        Response: cleanedBody(chirp.Chirp), 
    })    
    
}

func cleanedBody(msg string) string {

    filter := []string{"kerfuffle", "sharbert", "fornax"}
    body := strings.Split(msg, " ")
    for i, word := range body {
        for _, filter_word := range filter {
            if strings.ToLower(word) == filter_word {
                body[i] = "****"
            }
        }
    }
    
    return strings.Join(body, " ")
}

