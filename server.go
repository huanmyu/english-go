// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"git.oschina.net/bwn/english/chat"
	"git.oschina.net/bwn/english/model"
)

const timeLayout = "2006-01-02 15:04:05"

func WordsHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
	words, err := word.GetWordList(1, 5)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, words)
}

func WordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word := model.Word{ID: id}
	err = word.GetWordById(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, word)
}

func AddWordHandler(w http.ResponseWriter, r *http.Request) {
	var word Word
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&word); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	insertTime, err := time.Parse(timeLayout, "2017-03-01 13:00:00")
	if err != nil {
		log.Fatal(err)
	}
	// word := model.Word{
	// 	Name:        "peer",
	// 	Phonogram:   "[pɪr]",
	// 	Audio:       "https://tts.hjapi.com/en-us/9DF3BF531D7AD6B2",
	// 	Explanation: "n. 同等的人；同辈，同事 \nv. 凝视，盯着看",
	// 	Example:     "The WebSocket protocol defines three types of control messages: close, ping and pong. Call the connection WriteControl, WriteMessage or NextWriter methods to send a control message to the peer.",
	// 	CreatedAt:   insertTime,
	// 	UpdatedAt:   insertTime,
	// }
	err = word.CreateWord()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusCreated, word)
}

func EditWordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var word Word
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&w1); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	word.ID = id
	// word := model.Word{
	// 	ID:          3,
	// 	Name:        "wobble",
	// 	Phonogram:   "[ˈwɑbl]",
	// 	Explanation: "vt.& vi.<使>晃动； <使>摇摆不定； 颤动；\nn.摇动，晃动； 不稳定；",
	// 	Example:     "A topic that can easily make anyone's mind wobble.",
	// }
	err = word.UpdateWord()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, word)
}

func DeleteWordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/chat" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "static/chat.html")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithIndentJSON(w, code, map[string]string{"error": message})
}

func respondWithIndentJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		log.Fatalf("Json marshaling failed: %s", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Json marshaling failed: %s", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}

func main() {
	hub := chat.NewHub()
	go hub.Run()
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/words", WordsHandler)
	r.HandleFunc("/word/{id:[0-9]+}", WordHandler).Methods("GET")
	r.HandleFunc("/word", AddWordHandler).Methods("POST")
	r.HandleFunc("/word/{id:[0-9]+}", EditWordHandler).Methods("PUT")
	r.HandleFunc("/word/{id:[0-9]+}", DeleteWordHandler).Methods("DELETE")
	r.HandleFunc("/chat", chatHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
