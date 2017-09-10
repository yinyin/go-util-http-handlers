package httphandlers

import (
	"net/http"
	"time"
)

// Attach content modification related headers.
// Headers attached includes: ETag and Last-Modified.
func attachModificationTagHeader(w http.ResponseWriter, eTag string, modifyTime time.Time) {
	httpModifyTime := modifyTime.UTC().Format(http.TimeFormat)
	w.Header().Set("ETag", eTag)
	w.Header().Set("Last-Modified", httpModifyTime)
}

// Replies to the request with 304 Not Modified to indicate that the requested
// content is not modified since last request.
//
// The parameters eTag and modifyTime are ETag and modification time of content
// which will be attached to response headers.
func NotModify(w http.ResponseWriter, eTag string, modifyTime time.Time) {
	attachModificationTagHeader(w, eTag, modifyTime)
	http.Error(w, "304 Not Modified", http.StatusNotModified)
}

// Perform conditional-get on request r. The "If-None-Match" and "If-Modified-Since" headers
// will be extracted from request to compare with given ETag eTag and modification time modifyTime.
// A time.Duration acceptableAge is passed in for acceptable time difference on
// modification time comparision.
//
// The boolean true will be return if client content is updated and response is submitted.
// Otherwise, false will be return and the caller should serving the content to client.
func ConditionalGet(w http.ResponseWriter, r * http.Request, eTag string, modifyTime time.Time, acceptableAge time.Duration) (done bool) {
	remoteETag := r.Header.Get("If-None-Match")
	if remoteETag == eTag {
		NotModify(w, eTag, modifyTime)
		return true
	}
	remoteModifyTimeText := r.Header.Get("If-Modified-Since")
	if "" != remoteModifyTimeText {
		if remoteModifyTime, err := http.ParseTime(remoteModifyTimeText); nil == err {
			if modifyTime.Before(remoteModifyTime.Add(acceptableAge)) {
				NotModify(w, eTag, modifyTime)
				return true
			}
		}
	}
	return false
}
