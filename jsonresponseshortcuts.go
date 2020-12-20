package httphandlers

import (
	"net/http"
)

// JSONResponseWithStatusOK is shortcut of JSONResponseWithStatusCode(w, v, http.StatusOK) call.
func JSONResponseWithStatusOK(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusOK)
}

// JSONResponseWithStatusBadRequest is shortcut of JSONResponseWithStatusCode(w, v, http.StatusBadRequest) call.
func JSONResponseWithStatusBadRequest(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusBadRequest)
}

// JSONResponseWithStatusUnauthorized is shortcut of JSONResponseWithStatusCode(w, v, http.StatusUnauthorized) call.
func JSONResponseWithStatusUnauthorized(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusUnauthorized)
}

// JSONResponseWithStatusForbidden is shortcut of JSONResponseWithStatusCode(w, v, http.StatusForbidden) call.
func JSONResponseWithStatusForbidden(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusForbidden)
}

// JSONResponseWithStatusNotFound is shortcut of JSONResponseWithStatusCode(w, v, http.StatusNotFound) call.
func JSONResponseWithStatusNotFound(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusNotFound)
}

// JSONResponseWithStatusMethodNotAllowed is shortcut of JSONResponseWithStatusCode(w, v, http.StatusMethodNotAllowed) call.
func JSONResponseWithStatusMethodNotAllowed(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusMethodNotAllowed)
}

// JSONResponseWithStatusConflict is shortcut of JSONResponseWithStatusCode(w, v, http.StatusConflict) call.
func JSONResponseWithStatusConflict(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusConflict)
}

// JSONResponseWithStatusTooManyRequests is shortcut of JSONResponseWithStatusCode(w, v, http.StatusTooManyRequests) call.
func JSONResponseWithStatusTooManyRequests(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusTooManyRequests)
}

// JSONResponseWithStatusInternalServerError is shortcut of JSONResponseWithStatusCode(w, v, http.StatusInternalServerError) call.
func JSONResponseWithStatusInternalServerError(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusInternalServerError)
}

// JSONResponseWithStatusNotImplemented is shortcut of JSONResponseWithStatusCode(w, v, http.StatusNotImplemented) call.
func JSONResponseWithStatusNotImplemented(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusNotImplemented)
}

// JSONResponseWithStatusServiceUnavailable is shortcut of JSONResponseWithStatusCode(w, v, http.StatusServiceUnavailable) call.
func JSONResponseWithStatusServiceUnavailable(w http.ResponseWriter, v interface{}) (err error) {
	return JSONResponseWithStatusCode(w, v, http.StatusServiceUnavailable)
}
