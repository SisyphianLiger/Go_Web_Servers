package main

import (
	"encoding/json"
	"net/http"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/google/uuid"
)


func (cfg * apiConfig) upgradeToRed(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Event string `json:"event"` 
        Data struct {
            UserID uuid.UUID `json:"user_id"`
        } `json:"data"` 
    }
    
    _, apiError := auth.GetAPIKey(r.Header)
    if apiError != nil {
        respondWithError(w, http.StatusUnauthorized,"API Key is not correctly specified",apiError)
        return
    }

    params := parameters{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&params); err != nil {
        respondWithError(w, http.StatusNotFound,"Request Body Malformed",err)
        return
    }
    if params.Event != "user.upgraded" {
        respondWithJson(w, http.StatusNoContent, nil)
        return
    }

    _, uError := cfg.dbc.UpgradeUserToRed(r.Context(), params.Data.UserID)
    if uError != nil {
        respondWithError(w, http.StatusNotFound, "User not in DB", uError)
        return
    }
    
    respondWithJson(w, http.StatusNoContent, nil)

}
