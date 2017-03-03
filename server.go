// server.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"git.oschina.net/bwn/english/chat"
	"git.oschina.net/bwn/english/model"
)

const timeLayout = "2006-01-02 15:04:05"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
	// use request method
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pageNumberValue := r.FormValue("pageNumber")
	pageSizeValue := r.FormValue("pageSize")
	if pageNumberValue == "" {
		pageNumberValue = "1"
	}

	if pageSizeValue == "" {
		pageSizeValue = "20"
	}

	pageNumber, err := strconv.ParseInt(pageNumberValue, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeValue, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	words, err := word.GetWordList(pageNumber, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, words)
}

func wordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word := model.Word{ID: id}
	err = word.GetWordByID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, word)
}

func createWordHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&word); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	word.CreatedAt = time.Now()
	word.UpdatedAt = time.Now()
	lastInsertID, err := word.CreateWord()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	word.ID = lastInsertID
	respondWithIndentJSON(w, http.StatusCreated, word)
}

func editWordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}

	var word, viewWord model.Word
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&word); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	viewWord = model.Word{ID: id}
	err = viewWord.GetWordByID()
	if err != nil {
		return
	}

	word.CreatedAt = viewWord.CreatedAt
	word.UpdatedAt = time.Now()

	word.ID = id
	err = word.UpdateWord()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, word)
}

func deleteWordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid word ID")
		return
	}

	word := model.Word{ID: id}
	err = word.DeleteWord()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, map[string]string{"code": "200", "result": "success"})
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/chat.html")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithIndentJSON(w, code, map[string]string{"error": message})
}

func respondWithIndentJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}

func main() {
	hub := chat.NewHub()
	go hub.Run()
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/words", wordsHandler)
	r.HandleFunc("/word/{id:[0-9]+}", wordHandler).Methods("GET")
	r.HandleFunc("/word", createWordHandler).Methods("POST")
	r.HandleFunc("/word/{id:[0-9]+}", editWordHandler).Methods("PUT")
	r.HandleFunc("/word/{id:[0-9]+}", deleteWordHandler).Methods("DELETE")
	r.HandleFunc("/chat", chatHandler)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
