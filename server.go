// server.go
package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	//	headersCORSOption := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	//	originsCORSOption := handlers.AllowedOrigins([]string{"*"})
	//	methodsCORSOption := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS"})
	//  handlers.CORS(originsCORSOption, headersCORSOption, methodsCORSOption)(r),
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
