package v1

import (
	"net/http"
)

func (a *API) startHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the CRUD API!"))
}
