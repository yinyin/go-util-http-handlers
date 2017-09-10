package httphandlers

import (
	"net/http"
	"encoding/json"
	"strconv"
	"time"
)

func attachJSONContentHeader(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
}

func JSONResponse(w http.ResponseWriter, v interface{}) (err error) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "500 JSONResponse Failed:\n" + err.Error(), http.StatusInternalServerError)
		return err
	}
	attachJSONContentHeader(w, b)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return nil
}

func JSONResponseConditional(w http.ResponseWriter, r * http.Request, v interface{}, eTag string, modifyTime time.Time, acceptableAge time.Duration) (err error) {
	if ConditionalGet(w, r, eTag, modifyTime, acceptableAge) {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "500 JSONResponseConditional Failed:\n" + err.Error(), http.StatusInternalServerError)
		return err
	}
	attachModificationTagHeader(w, eTag, modifyTime)
	attachJSONContentHeader(w, b)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return nil
}