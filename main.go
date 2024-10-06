package main

import (
	"log"
	"net/http"
)


func main ()  {
    const filepathRoot = "."  
    const port = "8080"

    server:= http.NewServeMux()
    localSever := http.Server{
        Handler: server,
        Addr: ":" + port,
    }
    
    server.Handle("/",http.FileServer(http.Dir(filepathRoot)))

    log.Printf("Serving on port %s\n", port)
    log.Fatal(localSever.ListenAndServe())


}
