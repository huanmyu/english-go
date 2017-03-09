// chat.go
package controller

import (
	"net/http"

	"git.oschina.net/bwn/english/chat"
)

type Chat struct {
}

func (c Chat) ViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/chat.html")
}

func (c Chat) ServeHandler(w http.ResponseWriter, r *http.Request) {
	chat.ServeWs(hub, w, r)
}
