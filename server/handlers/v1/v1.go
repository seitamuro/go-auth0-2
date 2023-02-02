package v1

import (
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}
