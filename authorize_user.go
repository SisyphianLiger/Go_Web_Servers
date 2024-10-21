package main

import (
	"encoding/json"
	"net/http"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
)



func (cfg * apiConfig) authorizeUser(w http.ResponseWriter, r * http.Request) {
    type response struct {
        Password string `json:"password"`
        Email string `json:"email"`
    }

    
    token, RFerr := auth.GetBearerToken(r.Header)    
    if RFerr != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect Header no Refresh Token Found", RFerr)
        return
    }

    userID, jwtErr := auth.ValidateJWT(token, cfg.jwtSecret)
    if jwtErr != nil {
        respondWithError(w, http.StatusUnauthorized, "Unauthorized JwT Token", jwtErr)
        return
    }


    resp := response{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&resp); err != nil {
        respondWithError(w, http.StatusUnauthorized, "Malformed Body Information", err)
        return
    }

    pass, pErr := auth.HashPassword(resp.Password)
    if pErr != nil {

        respondWithError(w, http.StatusInternalServerError, "Unable to Hash password", pErr)
        return
    }
    
    dbUser := database.UpdateUserInfoParams{
        ID: userID,
        HashedPassword: pass,
        Email: resp.Email,
    }
    // DO something with Hashed Password?
    _, uErr := cfg.dbc.UpdateUserInfo(r.Context(), dbUser)
    if uErr != nil {
        respondWithError(w, http.StatusInternalServerError, "Unable to Update User email, check if user exists", pErr)
        return
    }

    // Respond 
    respondWithJson(w, http.StatusOK, resp)
}
