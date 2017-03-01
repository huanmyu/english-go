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

	"git.oschina.net/bwn/english/model"
)

const timeLayout = "2006-01-02 15:04:05"

func WordsHandler(w http.ResponseWriter, r *http.Request) {
	var word model.Word
	words := word.GetWordList(1, 5)
	data, err := json.MarshalIndent(words, "", "	")
	if err != nil {
		log.Fatalf("Json marshaling failed: %s", err)
	}
	fmt.Fprintf(w, "%s\n", data)
}

func WordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Fatalf("Strconv parseint failed: %s", err)
	}

	var word model.Word
	w1 := word.GetWordById(id)
	data, err := json.MarshalIndent(w1, "", "	")
	if err != nil {
		log.Fatalf("Json marshaling failed: %s", err)
	}
	fmt.Fprintf(w, "%s\n", data)
}

func AddWordHandler(w http.ResponseWriter, r *http.Request) {
	insertTime, err := time.Parse(timeLayout, "2017-03-01 13:00:00")
	if err != nil {
		log.Fatal(err)
	}
	word := model.Word{
		Name:        "tackle",
		Phonogram:   "[ˈtækəl] ",
		Audio:       "",
		Explanation: "n.用具，装备； 索具； 阻挡； 阻截队员 \nvt.着手处理； [橄榄球]擒住并摔倒（一名对方球员）； 给（马）配上挽具； \nvi.擒住并摔倒一名对手；",
		Example:     "Design patterns are solutions to recurring problems; guidelines on how to tackle certain problems.",
		CreatedAt:   insertTime,
		UpdatedAt:   insertTime,
	}
	lastID, rowCnt := word.CreateWord()
	fmt.Fprintf(w, "ID = %d, affected = %d\n", lastID, rowCnt)
}

func EditWordHandler(w http.ResponseWriter, r *http.Request) {
	word := model.Word{
		ID:          3,
		Name:        "wobble",
		Phonogram:   "[ˈwɑbl]",
		Explanation: "vt.& vi.<使>晃动； <使>摇摆不定； 颤动；\nn.摇动，晃动； 不稳定；",
		Example:     "A topic that can easily make anyone's mind wobble.",
	}
	rowCnt := word.UpdateWord()
	fmt.Fprintf(w, "ID = %d, affected = %d\n", word.ID, rowCnt)
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/words", WordsHandler)
	r.HandleFunc("/words/add", AddWordHandler)
	r.HandleFunc("/words/edit", EditWordHandler)
	r.HandleFunc("/words/{id:[0-9]+}", WordHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
