package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
        "time"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
	"github.com/google/uuid"
)


type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}


func (cfg * apiConfig) getAChirp(w http.ResponseWriter, r * http.Request) {
    chirpIDString := r.PathValue("chirpID")
    chirpID, err := uuid.Parse(chirpIDString)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
        return
    }
    // NO ROWS IN RESULT SET MEANS WE NEED TO CHECK IF THERE ARE ROWS AFTER ADDING
    chirp, err := cfg.dbc.GetAChirp(r.Context(), chirpID) 
    if err != nil {
        respondWithError(w, http.StatusNotFound, "DB Request Error", err)
        return 
    }

    respondWithJson(w, http.StatusOK, Chirp {
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
    })
}


func (cfg * apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
    chirps, err := cfg.dbc.GetChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "DB Request Error", err)
        return 
    }

    onlyChirps := []Chirp{}
   
    for _, dbChirp := range chirps {
        onlyChirps = append(onlyChirps, Chirp{
            ID: dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body: dbChirp.Body,
            UserID: dbChirp.UserID,
        })    
    }

    // Load Payload to chirp 
    respondWithJson(w, http.StatusOK, onlyChirps)
}

// Going to Redo this function with the types and stuff...
func (cfg * apiConfig) makeChirp (w http.ResponseWriter, r *http.Request){

    type parameters struct {
        Body   string `json:"body"`
        UserID uuid.UUID `json:"user_id"`
    }
    params := parameters{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&params); err != nil {
        respondWithError(w, http.StatusBadRequest, "Incorrect Json Object", err)
        return
    }

    cleanBody, err := cfg.validateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Body was empty", nil)
        return
    }
    
    chirp, err := cfg.dbc.MakeChirp(r.Context(), database.MakeChirpParams{Body: cleanBody, UserID: params.UserID}); 
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error(), err)
        return 
    }
    respondWithJson(w, http.StatusCreated, Chirp{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
    })    
}


func (cfg * apiConfig) validateChirp(body string) (string, error) {
    if len(body) > 140 {
        return "",fmt.Errorf("The Chirp is too long: MUST BE LESS THAN 140 Characters") 
    }
    return cleanedBody(body), nil

}

func cleanedBody(msg string) string {
    
    badWords := map[string]struct{}{
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }

    body := strings.Split(msg, " ")
    for i, word := range body {
        loweredWord := strings.ToLower(word)
        if _, ok := badWords[loweredWord]; ok {
            body[i] = "****"
        }
    }
    
    return strings.Join(body, " ")
}

