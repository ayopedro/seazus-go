package response

import (
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	WriteJSON(w, status, err.Error(), nil)
}
