package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/bowenchen6/english-go/model"
	"golang.org/x/net/publicsuffix"
)

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

func GetCookie() {
	// Start a server to give us cookies.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("Flavor"); err != nil {
			http.SetCookie(w, &http.Cookie{Name: "Flavor", Value: "Chocolate Chip"})
		} else {
			cookie.Value = "Oatmeal Raisin"
			http.SetCookie(w, cookie)
		}
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	if _, err = client.Get(u.String()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After 1st request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	if _, err = client.Get(u.String()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After 2nd request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
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
