package httphandlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type staticBytes struct {
	data               []byte
	contentType        string
	eTag               string
	modifyTime         time.Time
	acceptableAge      time.Duration
	httpModifyTimeText string
	lenText            string
}

// NewStaticBytes create http.Handler to serve static bytes with given
// `contentType`, `eTag`, `modifyTime` and `acceptableAge` for GET request.
//
// Given `data` will not be copy. The content of `data` should NOT be modified after
// invocation.
//
// If `data` is empty, the handler will response with 204 No Content.
//
// If `eTag` is empty string, it will be generated based on SHA-256 checksum of
// the data.
//
// If `modifyTime` is zero, it will be set to current time.
//
// Conditional GET and HEAD methods are supported.
func NewStaticBytes(data []byte, contentType, eTag string, modifyTime time.Time, acceptableAge time.Duration) (hnd http.Handler, err error) {
	if contentType == "" {
		err = ErrRequireContentType
		return
	}
	if eTag == "" {
		chksum := sha256.Sum256(data)
		eTag = `"` + base64.RawURLEncoding.EncodeToString(chksum[:]) + `"`
	} else if eTag[0] != '"' {
		eTag = `"` + eTag + `"`
	}
	if modifyTime.IsZero() {
		modifyTime = time.Now().UTC()
	}
	httpModifyTimeText := modifyTime.UTC().Format(http.TimeFormat)
	lenText := strconv.FormatInt(int64(len(data)), 10)
	hnd = &staticBytes{
		data:               data,
		contentType:        contentType,
		eTag:               eTag,
		modifyTime:         modifyTime,
		acceptableAge:      acceptableAge,
		httpModifyTimeText: httpModifyTimeText,
		lenText:            lenText,
	}
	return
}

func (hnd *staticBytes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(hnd.data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if done := ConditionalGet(
		w, r,
		hnd.eTag, hnd.modifyTime, hnd.acceptableAge); done {
		return
	}
	hdr := w.Header()
	hdr.Set("ETag", hnd.eTag)
	hdr.Set("Last-Modified", hnd.httpModifyTimeText)
	hdr.Set("Content-Type", hnd.contentType)
	hdr.Set("Content-Length", hnd.lenText)
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		w.Write(hnd.data)
	}
}

// NewStaticJSON create http.Handler to serve static JSON content `v` with given
// `eTag`, `modifyTime` and `acceptableAge` for GET and conditional GET request.
//
// The JSON is encoded with `encoding/json` package. If the encoding failed, an
// error will be returned.
//
// If `contentType` is empty, it will be set to "application/json; charset=utf-8".
//
// Result handler is backed by NewStaticBytes(). The characteristics are the same.
func NewStaticJSON(v interface{}, contentType, eTag string, modifyTime time.Time, acceptableAge time.Duration) (hnd http.Handler, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}
	return NewStaticBytes(b, contentType, eTag, modifyTime, acceptableAge)
}
