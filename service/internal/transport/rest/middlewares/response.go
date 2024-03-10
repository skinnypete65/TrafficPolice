package middlewares

import "net/http"

func authError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func notConfirmedError(w http.ResponseWriter) {
	http.Error(w, "You are not confirmed", http.StatusForbidden)
}
