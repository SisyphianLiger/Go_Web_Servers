package main

import (
    "fmt"
    "net/http"
    "log"
)


func (cfg * apiConfig) middlewareMetrics(next http.Handler) http.Handler {
    hitCounter := func (w http.ResponseWriter, r *http.Request){
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(hitCounter)
}


func (cfg * apiConfig) serverHits(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK) 

        body := fmt.Sprintf(`
        <html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>`, cfg.fileserverHits.Load())

        res, err := w.Write([]byte(body))
        if err != nil {
            log.Printf("Body not properly created: %d err: %v", res, err)
            return
        }
}
