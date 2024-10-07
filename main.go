package main

import (
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
    const apiPath = "/api"
    const adminPath = "/admin"
    // Config
    cfg := apiConfig{
        fileserverHits: atomic.Int32{},
    }
    
    // Server and Handlers
    server:= http.NewServeMux()
    server.Handle("/app/", cfg.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
    server.HandleFunc("GET "+ apiPath + "/healthz", healthCheck)
    server.HandleFunc("GET " + adminPath + "/metrics", cfg.serverHits)
    server.HandleFunc("POST " + adminPath + "/reset", cfg.serverReset)

    // Server Info
    localSever := http.Server{
        Handler: server,
        Addr: ":" + port,
    }
    
    log.Printf("Serving on port %s\n", port)
    log.Fatal(localSever.ListenAndServe())


}

