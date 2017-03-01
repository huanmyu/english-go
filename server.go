// server.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"git.oschina.net/bwn/english/model"
)

const timeLayout = "2006-01-02 15:04:05"

func handler(w http.ResponseWriter, r *http.Request) {
	//word := model.Word{}
	var word model.Word
	w1 := word.GetWordList()
	fmt.Fprintf(w, "Hi there, I love %s, at %s!", w1.Name, w1.CreatedAt.Format(timeLayout))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	insertTime, err := time.Parse(timeLayout, "2017-03-01 13:00:00")
	if err != nil {
		log.Fatal(err)
	}
	word := model.Word{
		Name:        "catchy",
		Phonogram:   "[ˈkætʃi, ˈkɛtʃi]",
		Audio:       "",
		Explanation: "adj.易记的； 易使人上当的； 迷人的；",
		Example:     "The songs were both catchy and original.",
		CreatedAt:   insertTime,
		UpdatedAt:   insertTime,
	}
	lastId, rowCnt := word.CreateWord()
	fmt.Fprintf(w, "ID = %d, affected = %d\n", lastId, rowCnt)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add", createHandler)
	http.ListenAndServe(":8080", nil)
}
