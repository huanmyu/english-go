package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bowenchen6/english-go/model"
	"github.com/gorilla/mux"
)

// User serve user handler
type User struct {
}

// LoginHandler serve user login
func (u User) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

// IsLoginHandler validate user is login
func (u User) IsLoginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := oauth(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, user)
}

// LogoutHandler user logout
func (u User) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("Token"); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		cookie.Expires = time.Now().AddDate(-1, 0, 0)
		http.SetCookie(w, cookie)
	}
	respondWithIndentJSON(w, http.StatusOK, map[string]string{"code": "200", "result": "success"})
}

// ListHandler find all user by page
func (u User) ListHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
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

	users, err := user.GetUserList(pageNumber, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithIndentJSON(w, http.StatusOK, users)
}

// ViewHandler find user by id
func (u User) ViewHandler(w http.ResponseWriter, r *http.Request) {
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

// CreateHandler create user
func (u User) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	if cookie, err := r.Cookie("Token"); err != nil {
		http.SetCookie(w, &http.Cookie{Name: "Token", Value: strconv.FormatInt(user.ID, 10)})
	} else {
		cookie.Value = strconv.FormatInt(user.ID, 10)
		http.SetCookie(w, cookie)
	}

	respondWithIndentJSON(w, http.StatusCreated, user)
}

// EditHandler edit user
func (u User) EditHandler(w http.ResponseWriter, r *http.Request) {
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
