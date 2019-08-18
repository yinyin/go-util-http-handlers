package httphandlers

import (
	"encoding/json"
	"net/http"
)

// DecodeJSONRequest parse JSON in the request body with given reference v.
// The HTTP error status code will be respond if decoding failed.
func DecodeJSONRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if nil == r.Body {
		http.Error(w, "empty request", http.StatusBadRequest)
		return ErrEmptyRequestBody
	}
	if err := json.NewDecoder(r.Body).Decode(v); nil != err {
		http.Error(w, "malformed request", http.StatusBadRequest)
		return err
	}
	return nil
}
