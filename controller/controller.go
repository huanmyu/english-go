package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"git.oschina.net/bwn/english/chat"

	"github.com/bowenchen6/english-go/model"
)

var hub *chat.Hub

func init() {
	hub = chat.NewHub()
	go hub.Run()
}

const salt = "word@english"

func oauth(r *http.Request) (user model.User, err error) {
	userId, err := r.Cookie("Token")
	if err != nil {
		return
	}

	user.ID, err = strconv.ParseInt(userId.Value, 10, 64)
	if err != nil {
		return
	}

	err = user.GetUserByID()
	return
}

func md5Password(password string) string {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(string(salt)))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
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
