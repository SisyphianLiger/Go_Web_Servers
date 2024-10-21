package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
)


func (cfg * apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
    
    result, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Token in DB not found", err)
        return
    }

    refToken, error := cfg.dbc.GetRefreshToken(r.Context(), result)
    if error != nil {
        respondWithError(w, http.StatusInternalServerError, "Refresh Token Not Found", err)
        return
    }
    databaseResponse := database.UpdateRefreshTokenParams{
        RevokedAt: sql.NullTime{
            Time: time.Now(),
            Valid: true,
        },
        UpdatedAt: time.Now(),
        Token: refToken.Token,
    }

    // Post to DB
    _, failure := cfg.dbc.UpdateRefreshToken(r.Context(), databaseResponse) 
    if failure != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to Update DB", err)
        return
    }
        
    respondWithJson(w, http.StatusNoContent, nil)

}
