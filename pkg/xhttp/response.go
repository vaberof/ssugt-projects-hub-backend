package xhttp

import (
	"encoding/json"
	"net/http"
)

func WriteResponseJson(w http.ResponseWriter, status int, a any) error {
	if a == nil {
		return nil
	}

	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set(AccessControlAllowOrigin, "*")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}
