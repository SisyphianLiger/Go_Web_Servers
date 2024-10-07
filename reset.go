package main

import (
    "net/http"
    "fmt"
    "log"
)

func (cfg * apiConfig) serverReset(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK) 
        cfg.fileserverHits.CompareAndSwap(cfg.fileserverHits.Load(), 0)
        res, err := w.Write([]byte(fmt.Sprintf("Server Hits Set Back to %d", cfg.fileserverHits.Load()))) 
        if err != nil {
        log.Printf("Body not properly created: %d err: %v", res, err)
            return
        }
}

