package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
	"github.com/google/uuid"
)

type refreshtoken struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
    Token string `json:"token"`
    RefreshToken string `json:"refresh_token"`
    Is_Chirpy_Red bool `json:"is_chirpy_red"`
}

func (cfg * apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {
    // Process Email Login First
    type parameters struct {
        Email string
        Password string
        ExpiresInSeconds *time.Duration   `json:"expires_in_seconds,omitempty"`
    }
    type response struct {
        refreshtoken
    }

    p  := parameters{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Inccorrect Body", err)
        return
    }


    // Make Query for sqlc (Find User by Email)
    user, uErr := cfg.dbc.GetUser(r.Context(), p.Email)
    if uErr != nil {
        respondWithError(w, http.StatusInternalServerError, uErr.Error(), uErr)
        return
    }
    
    if invalidLoginErr := auth.CheckPasswordHash(user.HashedPassword, p.Password); invalidLoginErr != nil {
        respondWithError(w, http.StatusUnauthorized, "Not a valid password, please try again", invalidLoginErr)
        return
    }

    token, terr := auth.MakeJWT(user.ID, cfg.jwtSecret)
    if terr != nil {
        respondWithError(w, http.StatusInternalServerError, "Error generating token", terr)
        return
    }


    refToken, rerr := auth.MakeRefreshToken()
    if rerr != nil {
        respondWithError(w, http.StatusInternalServerError, "Error generating token", rerr)
        return
    }

    databaseToken := database.CreateRefreshTokenParams{
        Token: refToken,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID,
        ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
        RevokedAt: sql.NullTime{
            Time: time.Time{},
            Valid: false,
        },
    }
    // Gonna add it here
    _, TokenError := cfg.dbc.CreateRefreshToken(r.Context(), databaseToken)
    if TokenError != nil {
        respondWithError(w, http.StatusInternalServerError, "Error Adding Token to DB", TokenError)
        return
    }

    respondWithJson(w, http.StatusOK, response{
        refreshtoken {
            ID: user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: user.Email,
            Token: token,
            RefreshToken: refToken,
            Is_Chirpy_Red: user.IsChirpyRed.Bool,
        },
    })
}
