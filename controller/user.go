package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/bowenchen6/english/model"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user.Password = md5Password(user.Password)
	err := user.GetUserByNameAndPassword()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if cookie, err := r.Cookie("Token"); err != nil {
		http.SetCookie(w, &http.Cookie{Name: "Token", Value: strconv.FormatInt(user.ID, 10)})
	} else {
		cookie.Value = strconv.FormatInt(user.ID, 10)
		http.SetCookie(w, cookie)
	}

	respondWithIndentJSON(w, http.StatusOK, user)

}

// UserHandler find user by id
func UserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user := model.User{ID: id}
	err = user.GetUserByID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, user)
}

// CreateUserHandler create user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user.Password = md5Password(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	lastInsertID, err := user.CreateUser()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.ID = lastInsertID
	respondWithIndentJSON(w, http.StatusCreated, user)
}

// EditUserHandler edit user
func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user, viewuser model.User
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	viewuser = model.User{ID: id}
	err = viewuser.GetUserByID()
	if err != nil {
		return
	}

	user.CreatedAt = viewuser.CreatedAt
	user.UpdatedAt = time.Now()

	user.ID = id
	err = user.UpdateUser()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, user)
}
