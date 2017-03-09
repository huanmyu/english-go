// router
package main

import (
	"github.com/bowenchen6/english-go/controller"
	"github.com/gorilla/mux"
)

var r, s *mux.Router

func init() {
	r = mux.NewRouter()

	r.HandleFunc("/", controller.HomeHandler)
	r.HandleFunc("/words", controller.WordListHandler)
	//	r.HandleFunc("/chat", chat.ViewHandler)
	//	r.HandleFunc("/ws", chat.ServeHandler)

	//	s = r.PathPrefix("/api").Subrouter().StrictSlash(true)
	//	userRouter()
	//	wordRouter()
}

//func userRouter() {
//	s.HandleFunc("/login", user.LoginHandler).Methods("POST")
//	s.HandleFunc("/isLogin", user.IsLoginHandler).Methods("GET")
//	s.HandleFunc("/user", user.CreateHandler).Methods("POST")
//	s.HandleFunc("/user/{id:[0-9]+}", user.ViewHandler).Methods("GET")
//	s.HandleFunc("/user/{id:[0-9]+}", user.EditHandler).Methods("PUT")
//}

//func wordRouter() {
//	s.HandleFunc("/words", word.ListHandler)
//	s.HandleFunc("/word/{id:[0-9]+}", word.ViewHandler).Methods("GET")
//	s.HandleFunc("/word", word.CreateHandler).Methods("POST")
//	s.HandleFunc("/word/{id:[0-9]+}", word.EditHandler).Methods("PUT")
//	s.HandleFunc("/word/{id:[0-9]+}", word.DeleteHandler).Methods("DELETE")
//}
