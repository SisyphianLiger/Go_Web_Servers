package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// apiConfig Struct
type apiConfig struct {
    fileserverHits atomic.Int32
}



func main ()  {
    // Important FilePaths/Addresses
    const filepathRoot = "."  
    const port = "8080"

    // Config
    cfg := apiConfig{
        fileserverHits: atomic.Int32{},
    }
    
    // Server and Handlers
    server:= http.NewServeMux()
    server.Handle("/app/", cfg.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
    server.HandleFunc("/healthz", healthCheck)
    server.HandleFunc("/metrics", cfg.serverHits)
    server.HandleFunc("/reset", cfg.serverReset)

    // Server Info
    localSever := http.Server{
        Handler: server,
        Addr: ":" + port,
    }
    
    log.Printf("Serving on port %s\n", port)
    log.Fatal(localSever.ListenAndServe())


}

func (cfg * apiConfig) middlewareMetrics(next http.Handler) http.Handler {
    hitCounter := func (w http.ResponseWriter, r *http.Request){
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(hitCounter)
}

func (cfg * apiConfig) serverHits(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK) 
        res, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
        if err != nil {
        log.Printf("Body not properly created: %d err: %v", res, err)
            return
        }
}
