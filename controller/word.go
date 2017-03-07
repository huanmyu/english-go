package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bowenchen6/english/model"
	"github.com/gorilla/mux"
)

// WordsHandler find all word by page
func WordsHandler(w http.ResponseWriter, r *http.Request) {
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

// WordHandler find word by id
func WordHandler(w http.ResponseWriter, r *http.Request) {
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

// CreateWordHandler create word
func CreateWordHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

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

// EditWordHandler edit word
func EditWordHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

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

// DeleteWordHandler delete word
func DeleteWordHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

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
