package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
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

type MakeChirpParams struct {
    Body   string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}
// Going to Redo this function with the types and stuff...
func (cfg * apiConfig) makeChirp (w http.ResponseWriter, r *http.Request){

    params := MakeChirpParams{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&params); err != nil {
        respondWithError(w, http.StatusBadRequest, "Expecting a ID and Body got incorrect Object", err)
        return
    }

    res,_ := auth.GetBearerToken(r.Header)
    log.Printf("Ummm what is the token in the header %s", res)
    // CHECK HERE
    // Helper Function to clean Function
    userID, verificationError := cfg.verifyJWT(r)
    if verificationError != nil {
        respondWithError(w, http.StatusUnauthorized, "JWT Token has expired:", verificationError)
        return
    }

    cleanBody, err := cfg.validateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Body was empty", nil)
        return
    }
    
    chirp, err := cfg.dbc.MakeChirp(r.Context(), database.MakeChirpParams{Body: cleanBody, UserID: userID}); 
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


func (cfg * apiConfig) verifyJWT(r *http.Request) (uuid.UUID, error) {

    validToken, err := auth.GetBearerToken(r.Header)
    if err != nil {
        return uuid.Nil, err
    }

    id, invalidID := auth.ValidateJWT(validToken, cfg.jwtSecret)
    if invalidID != nil {
        return uuid.Nil,invalidID
    }

    return id, nil
}
