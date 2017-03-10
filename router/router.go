package router

import (
	"github.com/bowenchen6/english-go/controller"
	"github.com/gorilla/mux"
)

// R is server router
var R, s *mux.Router

func init() {
	R = mux.NewRouter()

	homeRouter()
	chatRouter()

	s = R.PathPrefix("/api").Subrouter().StrictSlash(true)

	userRouter()
	wordRouter()
}

func homeRouter() {
	var h controller.Home
	R.HandleFunc("/", h.ViewHandler)
}

func chatRouter() {
	var c controller.Chat
	R.HandleFunc("/chat", c.ViewHandler)
	R.HandleFunc("/ws", c.ServeHandler)
}

func userRouter() {
	var u controller.User
	s.HandleFunc("/login", u.LoginHandler).Methods("POST")
	s.HandleFunc("/isLogin", u.IsLoginHandler).Methods("GET")
	s.HandleFunc("/users", u.ListHandler).Methods("GET")
	s.HandleFunc("/user", u.CreateHandler).Methods("POST")
	s.HandleFunc("/user/{id:[0-9]+}", u.ViewHandler).Methods("GET")
	s.HandleFunc("/user/{id:[0-9]+}", u.EditHandler).Methods("PUT")
}

func wordRouter() {
	var w controller.Word
	s.HandleFunc("/words", w.ListHandler).Methods("GET")
	s.HandleFunc("/word/{id:[0-9]+}", w.ViewHandler).Methods("GET")
	s.HandleFunc("/word", w.CreateHandler).Methods("POST")
	s.HandleFunc("/word/{id:[0-9]+}", w.EditHandler).Methods("PUT")
	s.HandleFunc("/word/{id:[0-9]+}", w.DeleteHandler).Methods("DELETE")
}
