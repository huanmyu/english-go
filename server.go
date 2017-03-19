// server.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bowenchen6/english-go/model"
	"github.com/bowenchen6/english-go/router"
)

func main() {
	defer model.R.Close()
	defer model.DB.Close()
	//	headersCORSOption := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	//	originsCORSOption := handlers.AllowedOrigins([]string{"*"})
	//	methodsCORSOption := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS"})
	//  handlers.CORS(originsCORSOption, headersCORSOption, methodsCORSOption)(r),
	srv := &http.Server{
		Handler:      router.N,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
