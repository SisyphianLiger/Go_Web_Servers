package main

import (
	"github.com/SisyphianLiger/Go_Web_Servers/internal/database"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync/atomic"
)

// apiConfig Struct
type apiConfig struct {
	fileserverHits     atomic.Int32
	dbc *database.Queries
	dev string
}


func main() {
	// Important FilePaths/Addresses
	const filepathRoot = "."
	const port = "8080"
	const apiPath = "/api"
	const adminPath = "/admin"

	openEnv()
	dbURL := environmentVarExists("DB_URL")
	devEnv := environmentVarExists("PLATFORM")

	// Make DB Connection extracted
	db := openDB("postgres", dbURL)
	dbQueries := database.New(db)

	// Config
	cfg := apiConfig{
		fileserverHits:     atomic.Int32{},
		dbc: dbQueries,
		dev: devEnv,
	}

	// Server and Handlers
	server := http.NewServeMux()
	server.Handle("/app/", cfg.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	server.HandleFunc("GET "+apiPath+"/healthz", healthCheck)
	server.HandleFunc("GET "+apiPath+"/chirps", cfg.getChirps)
	server.HandleFunc("GET "+apiPath+"/chirps/{chirpID}", cfg.getAChirp)
	server.HandleFunc("POST "+apiPath+"/chirps", cfg.makeChirp)
	server.HandleFunc("POST "+apiPath+"/users", cfg.addUser)
	server.HandleFunc("POST "+adminPath+"/reset", cfg.resetUserTable)
	server.HandleFunc("GET "+adminPath+"/metrics", cfg.serverHits)
	// server.HandleFunc("POST "+adminPath+"/reset", cfg.serverReset)

	// Server Info
	localSever := http.Server{
		Handler: server,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(localSever.ListenAndServe())
}
