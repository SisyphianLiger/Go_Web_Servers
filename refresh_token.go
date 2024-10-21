package main

import (
	"net/http"
	"time"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
)


func (cfg * apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
    type response struct {
        Token string `json:"token"`
    }

    currentTime := time.Now()
    
    res, RFerr := auth.GetBearerToken(r.Header)    
    if RFerr != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect Header no Refresh Token Found", RFerr)
        return
    }

    getToken, TErr := cfg.dbc.GetRefreshToken(r.Context(), res)
    if TErr != nil {
        respondWithError(w, http.StatusUnauthorized, "Unable to find Refresh Token in Database", TErr)
        return
    }
   
    if currentTime.After(getToken.ExpiresAt) {
        respondWithError(w, http.StatusUnauthorized, "Token Has Expired", TErr)
        return
    }

    if getToken.RevokedAt.Valid {
        respondWithError(w, http.StatusUnauthorized, "Token Has Been Revoked", nil)
        return
    }
    // After validating the refresh token...
    userID := getToken.UserID
    newAccessToken, err := auth.MakeJWT(userID, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Unable to Create new AccessToken", err)
        return
    }
    // HERE WE NEED TO GENERATE NEW ACCESS TOKEN
    // CHECK IF TOKEN IS EXPIRED?
    respondWithJson(w, http.StatusOK, response{
        Token: newAccessToken,
    })
}
