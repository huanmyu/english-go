package router

import (
	"github.com/bowenchen6/english-go/controller"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// N is server handler router and middlewares
var N *negroni.Negroni
var r, s *mux.Router

func init() {
	r = mux.NewRouter()
	homeRouter()
	chatRouter()

	s = r.PathPrefix("/api").Subrouter().StrictSlash(true)
	userRouter()
	wordRouter()

	// Includes some default(logger, recovery and static) middlewares
	N = negroni.Classic()
	N.UseHandler(r)
}

func homeRouter() {
	var h controller.Home
	r.HandleFunc("/", h.ViewHandler)
}

func chatRouter() {
	var c controller.Chat
	r.HandleFunc("/chat", c.ViewHandler)
	r.HandleFunc("/ws", c.ServeHandler)
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
