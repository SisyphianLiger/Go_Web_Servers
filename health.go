package main
import (
    "net/http"
    "log"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK) 
        res, err := w.Write([]byte(http.StatusText(http.StatusOK))) 
        if err != nil {
        log.Printf("Body not properly created: %d err: %v", res, err)
            return
        }
}

