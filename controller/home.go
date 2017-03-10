package controller

import "net/http"

// Home serve home handler
type Home struct {
}

// ViewHandler show home page
func (h Home) ViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}
