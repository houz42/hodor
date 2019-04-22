package handler

import (
	"net/http"
	"strings"
)

func setTempCookie(w http.ResponseWriter, r *http.Request, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  key,
		Value: value,
		// Domain:   domainName(r),
		MaxAge: stateTTL,
		Secure: true,
		// HttpOnly: true,
	})
}

func domainName(r *http.Request) string {
	host := r.Host
	if idx := strings.IndexByte(host, ':'); idx > 0 {
		host = host[:idx]
	}
	return host
}
