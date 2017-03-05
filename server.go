// server.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"git.oschina.net/bwn/english/chat"
	"git.oschina.net/bwn/english/controller"
)

const timeLayout = "2006-01-02 15:04:05"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/chat.html")
}

func main() {
	hub := chat.NewHub()
	go hub.Run()
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/user", controller.CreateUserHandler).Methods("POST")
	r.HandleFunc("/user/{id:[0-9]+}", controller.UserHandler).Methods("GET")
	r.HandleFunc("/user/{id:[0-9]+}", controller.EditUserHandler).Methods("PUT")

	r.HandleFunc("/words", controller.WordsHandler)
	r.HandleFunc("/word/{id:[0-9]+}", controller.WordHandler).Methods("GET")
	r.HandleFunc("/word", controller.CreateWordHandler).Methods("POST")
	r.HandleFunc("/word/{id:[0-9]+}", controller.EditWordHandler).Methods("PUT")
	r.HandleFunc("/word/{id:[0-9]+}", controller.DeleteWordHandler).Methods("DELETE")

	r.HandleFunc("/chat", chatHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	headersCORSOption := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsCORSOption := handlers.AllowedOrigins([]string{"*"})
	methodsCORSOption := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS"})

	srv := &http.Server{
		Handler:      handlers.CORS(originsCORSOption, headersCORSOption, methodsCORSOption)(r),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
