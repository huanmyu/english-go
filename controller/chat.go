package controller

import (
	"net/http"

	"github.com/bowenchen6/english-go/chat"
)

// Chat server websocket handler
type Chat struct {
}

// ViewHandler show chat page
func (c Chat) ViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/chat.html")
}

// ServeHandler serve websocket
func (c Chat) ServeHandler(w http.ResponseWriter, r *http.Request) {
	chat.ServeWs(hub, w, r)
}
