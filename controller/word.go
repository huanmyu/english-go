package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bowenchen6/english-go/model"
	"github.com/gorilla/mux"
)

// Word serve word handler
type Word struct {
}

// LatestListHandler show latest created words list
func (wd Word) LatestListHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
	var userID string

	if cookie, err := r.Cookie("Token"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		userID = cookie.Value
	}

	words, err := word.GetLatestCreatedWords(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, words)
}

// ListHandler find all word by page
func (wd Word) ListHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
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

// ViewHandler find word by id
func (wd Word) ViewHandler(w http.ResponseWriter, r *http.Request) {
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

// CreateHandler create word
func (wd Word) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
	var userID string

	if cookie, err := r.Cookie("Token"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		userID = cookie.Value
	}

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

	// cache latest created word for user
	wordID := strconv.FormatInt(lastInsertID, 10)
	err = word.CacheLatestCreatedWords(userID, wordID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	word.ID = lastInsertID
	respondWithIndentJSON(w, http.StatusCreated, word)
}

// EditHandler edit word
func (wd Word) EditHandler(w http.ResponseWriter, r *http.Request) {
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

// DeleteHandler delete word
func (wd Word) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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
