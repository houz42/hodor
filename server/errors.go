package server

import "net/http"

func internalError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func unauthorizedError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(msg))
}
