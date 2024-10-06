package main

import (
	"fmt"
	"net/http"
)


func main ()  {
    server:= http.NewServeMux()
    localSever := http.Server{
        Handler: server,
        Addr: ":8080",
    }
    err := localSever.ListenAndServe()
    if err != nil {
        e := fmt.Errorf("Error Starting Sever %s", err) 
        fmt.Println(e.Error())
    }

}
