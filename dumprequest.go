package httphandlers

import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
)

// Wrapping of https://godoc.org/net/url#Userinfo
type userinfoDataWrap struct {
	Username string
	Password string
	PasswordSet bool
	StringRepresentation string
}

func newUserinfoDataWrap(d * url.Userinfo) (r * userinfoDataWrap) {
	if nil == d {
		return nil
	}
	passwordText, passwordIsSet := d.Password()
	return &userinfoDataWrap {
		Username: d.Username(),
		Password: passwordText,
		PasswordSet: passwordIsSet,
		StringRepresentation:d.String(),
	}
}

// Wrapping of https://godoc.org/net/url#URL
type urlDataWrap struct {
	Scheme     string
	Opaque     string
	User       *userinfoDataWrap
	Host       string
	Path       string
	RawPath    string
	ForceQuery bool
	RawQuery   string
	Fragment   string
}

func newURLDataWrap(d * url.URL) (r * urlDataWrap) {
	if nil == d {
		return nil
	}
	return &urlDataWrap {
		Scheme: d.Scheme,
		Opaque: d.Opaque,
		User: newUserinfoDataWrap(d.User),
		Host       : d.Host,
		Path       : d.Path,
		RawPath: d.RawPath,
		ForceQuery: d.ForceQuery,
		RawQuery: d.RawQuery,
		Fragment: d.Fragment,
	}
}

// Wrapping of https://godoc.org/net/http#Request
type requestDataWrap struct {
	Method string
	URL *urlDataWrap
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     map[string][]string
	ContentLength int64
	TransferEncoding []string
	Close bool
	Host string
	RemoteAddr string
	RequestURI string
	HasTLSConnectionState bool
}

func newRequestDataWrap(d * http.Request) (r * requestDataWrap) {
	hasTLSConnectionState := false
	if nil != d.TLS {
		hasTLSConnectionState = true
	}
	return & requestDataWrap{
		Method:           d.Method,
		URL:              newURLDataWrap(d.URL),
		Proto:            d.Proto,
		ProtoMajor:       d.ProtoMajor,
		ProtoMinor:       d.ProtoMinor,
		Header:           d.Header,
		ContentLength:    d.ContentLength,
		TransferEncoding: d.TransferEncoding,
		Close:            d.Close,
		Host:             d.Host,
		RemoteAddr: d.RemoteAddr,
		RequestURI: d.RequestURI,
		HasTLSConnectionState: hasTLSConnectionState,
	}
}

type dumpRequestHandler struct {
}

// Create a handler which dumps request to response. The result will be text
// or JSON depends on request path.
//
// The response will be plain text in default. If the first part of path
// is `/json` (eg: /json/test...) the response will given in JSON.
func NewDumpRequestHandler() (h http.Handler) {
	return &dumpRequestHandler {}
}

func (h * dumpRequestHandler) serveJSON(w http.ResponseWriter, r * http.Request) {
	JSONResponse(w, newRequestDataWrap(r))
}

func (h * dumpRequestHandler) serveText(w http.ResponseWriter, r * http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w,"Method: %v\n", r.Method)
	fmt.Fprintf(w,"URL: (%v)\n", r.URL)
	fmt.Fprintf(w,"    Scheme: %v\n", r.URL.Scheme)
	fmt.Fprintf(w,"    Opaque: %v\n", r.URL.Opaque)
	fmt.Fprintf(w,"    User: %v\n", r.URL.User)
	fmt.Fprintf(w,"    Host: %v\n", r.URL.Host)
	fmt.Fprintf(w,"    Path: %v\n", r.URL.Path)
	fmt.Fprintf(w,"    RawPath: %v\n", r.URL.RawPath)
	fmt.Fprintf(w,"    ForceQuery: %v\n", r.URL.ForceQuery)
	fmt.Fprintf(w,"    RawQuery: %v\n", r.URL.RawQuery)
	fmt.Fprintf(w,"    Fragment: %v\n", r.URL.Fragment)
	fmt.Fprintf(w,"Proto: %v (major: %v, minor: %v)\n", r.Proto, r.ProtoMajor, r.ProtoMinor)
	fmt.Fprintf(w,"Header: (%v items)\n", len(r.Header))
	for headerKey, headerValue := range r.Header {
		fmt.Fprintf(w,"    %v: %v\n", headerKey, headerValue)
	}
	fmt.Fprintf(w,"ContentLength: %v\n", r.ContentLength)
	fmt.Fprintf(w,"TransferEncoding: %v\n", r.TransferEncoding)
	fmt.Fprintf(w,"Close: %v\n", r.Close)
	fmt.Fprintf(w,"Host: %v\n", r.Host)
	fmt.Fprintf(w,"RemoteAddr: %v\n", r.RemoteAddr)
	fmt.Fprintf(w,"RequestURI: %v\n", r.RequestURI)
	fmt.Fprintf(w,"TLS: %v\n", r.TLS)
}

func (h * dumpRequestHandler) ServeHTTP(w http.ResponseWriter, r * http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/json") {
		h.serveJSON(w, r)
	} else {
		h.serveText(w, r)
	}
}
