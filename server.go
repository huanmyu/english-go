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
	var word model.Word
	w1 := word.GetWordById(1)
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
	lastID, rowCnt := word.CreateWord()
	fmt.Fprintf(w, "ID = %d, affected = %d\n", lastID, rowCnt)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	word := model.Word{
		ID:          5,
		Name:        "tackle",
		Phonogram:   "[ˈtækəl]",
		Explanation: "n.用具，装备； 索具； 阻挡； 阻截队员 \nvt.着手处理； [橄榄球]擒住并摔倒（一名对方球员）； 给（马）配上挽具；\nvi.擒住并摔倒一名对手；",
		Example:     "Design patterns are solutions to recurring problems; guidelines on how to tackle certain problems.",
	}
	rowCnt := word.UpdateWord()
	fmt.Fprintf(w, "ID = %d, affected = %d\n", word.ID, rowCnt)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add", createHandler)
	http.HandleFunc("/edit", updateHandler)
	http.ListenAndServe(":8080", nil)
}
