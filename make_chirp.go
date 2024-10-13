package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
    Body   string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func (cfg * apiConfig) makeChirp(w http.ResponseWriter, r *http.Request){

    chirp := cfg.validateChirp(w, r)
    if chirp == (Chirp{}) {
        respondWithError(w, http.StatusBadRequest, "Body was empty", nil)
        return
    }
    
    err := cfg.dbc.MakeChirp(r.Context(), database.MakeChirpParams{Body: chirp.Body, UserID: chirp.UserID}); 
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "DB Request Error", err)
        return 
    }
    respondWithJson(w, http.StatusCreated, chirp)    
}


func (cfg * apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) Chirp {
    chirp := Chirp{} 

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&chirp); err != nil {
        respondWithError(w, http.StatusBadRequest, "Body of Request not acceptable Parameters", err)
        return Chirp{} 
    }

    if len(chirp.Body) > 140 {
        respondWithError(w, http.StatusBadRequest,"Chirp is too long", nil) 
        return Chirp{}
    }

    cleanedBody(chirp.Body)
    return Chirp{
        Body: chirp.Body, 
        UserID: chirp.UserID,
    }

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

