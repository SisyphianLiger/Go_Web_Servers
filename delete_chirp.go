package main

import (
	"net/http"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/auth"
	"github.com/google/uuid"
)



func (cfg * apiConfig) deleteAChirp(w http.ResponseWriter, r *http.Request) {
    
    chirpIDString := r.PathValue("chirpID")
    chirpID, err := uuid.Parse(chirpIDString)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "ID is not correct", err)
    }

    authToken, aErr := auth.GetBearerToken(r.Header)        
    if aErr != nil {
        respondWithError(w, http.StatusUnauthorized, "Malformed Token", aErr)
        return
    }
    
    userID, uErr := auth.ValidateJWT(authToken, cfg.jwtSecret)
    if uErr != nil {
        respondWithError(w, http.StatusForbidden, "Invalid Token Found", uErr)
        return
    }

    // We need to get the chirp and reference it to userID
    cID, cErr := cfg.dbc.GetAChirp(r.Context(), chirpID)
    if cErr != nil {
        respondWithError(w, http.StatusBadRequest, "Chirp could not be found", cErr)
        return
    }

    if cID.UserID != userID {
        respondWithError(w, http.StatusForbidden, "User is not authorized to delete tweet", nil)
        return
    }


    if deleteError := cfg.dbc.DeleteChirp(r.Context(), chirpID); deleteError != nil {
        respondWithError(w, http.StatusNotFound, "No Chirp with specified Name", deleteError)
        return
    }

    respondWithJson(w, http.StatusNoContent, nil)


}
