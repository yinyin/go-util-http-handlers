package httphandlers

import (
	"net/http"
	"strconv"
	"time"
)

func attachUTF8HTMLContentHeader(w http.ResponseWriter, htmlContent []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(htmlContent)), 10))
}

// UTF8HTMLResponseWithStatusCode generate HTML response based on given `htmlContent` which
// must encoded with UTF-8 charset. The no-cache headers enabled.
func UTF8HTMLResponseWithStatusCode(w http.ResponseWriter, htmlContent []byte, statusCode int) (err error) {
	attachUTF8HTMLContentHeader(w, htmlContent)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(statusCode)
	w.Write(htmlContent)
	return nil
}

// UTF8HTMLResponseConditional generate HTML response based on given `htmlContent` and handle
// conditional GET request. The given `htmlContent` must encoded with UTF-8 charset.
func UTF8HTMLResponseConditional(w http.ResponseWriter, r *http.Request, htmlContent []byte, eTag string, modifyTime time.Time, acceptableAge time.Duration) (err error) {
	if ConditionalGet(w, r, eTag, modifyTime, acceptableAge) {
		return nil
	}
	attachModificationTagHeader(w, eTag, modifyTime)
	attachUTF8HTMLContentHeader(w, htmlContent)
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		w.Write(htmlContent)
	}
	return nil
}
