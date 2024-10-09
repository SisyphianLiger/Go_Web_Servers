package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
	_ "github.com/lib/pq"
)

// apiConfig Struct
type apiConfig struct {
    fileserverHits atomic.Int32
    databaseConnection *database.Queries
}



func main ()  {
    // Important FilePaths/Addresses
    const filepathRoot = "."  
    const port = "8080"
    const apiPath = "/api"
    const adminPath = "/admin"

    // DB File
    dbURL := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Printf("Unable to Open DB, check ports and connecion string")
        return
    }
    dbQueries := database.New(db)

    // Config
    cfg := apiConfig{
        fileserverHits: atomic.Int32{},
        databaseConnection: dbQueries,
    }
    
    // Server and Handlers
    server:= http.NewServeMux()
    server.Handle("/app/", cfg.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
    server.HandleFunc("GET "+ apiPath + "/healthz", healthCheck)
    server.HandleFunc("POST "+ apiPath + "/validate_chirp", validateChirp)
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

