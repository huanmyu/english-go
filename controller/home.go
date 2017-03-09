package controller

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Getwd()
	log.Println("current path:", file)
	var buffer bytes.Buffer
	buffer.WriteString(file)
	buffer.WriteString("static/home.html")
	p := buffer.String()
	//respondWithError(w, 200, os.Args[0])
	http.ServeFile(w, r, p)
}
