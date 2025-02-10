package xhttp

import (
	"encoding/json"
	"net/http"
)

func ReadRequestJson(r *http.Request, a any) error {
	if a == nil {
		return nil
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(a)
}
