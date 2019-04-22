package handler

import (
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

type Config struct {
	Rules        map[string]*OAuthRule `json:"rules,omitempty"`
	VerifyPath   string                `json:"verify_path,omitempty"`
	CallbackPath string                `json:"callback_path,omitempty"`
}

func Handler(cfg *Config, l logr.Logger) http.Handler {
	for _, r := range cfg.Rules {
		r.log = l
	}

	mux := http.NewServeMux()

	mux.HandleFunc(cfg.VerifyPath, func(w http.ResponseWriter, r *http.Request) {
		log := l.WithValues("host", r.Host, "path", r.URL.Path)
		rule := cfg.Rules[r.Host]
		if rule == nil {
			log.V(2).Info("no rule")
			unauthorized(w, "no rule configured")
			return
		}

		rule.Verify(w, r)
	})

	mux.HandleFunc(cfg.CallbackPath, func(w http.ResponseWriter, r *http.Request) {
		log := l.WithValues("host", r.Host, "path", r.URL.Path)
		log.V(5).Info("request", "cookies", fmt.Sprintf("%v", r.Cookies()))
		hc, e := r.Cookie(hostCookieKey)
		if e != nil || hc.Value == "" {
			log.Error(e, "get original host cookie failed")
			badRequest(w, "unknown original host")
			return
		}
		host := hc.Value
		log.WithValues("origin host", host)

		rule := cfg.Rules[host]
		if rule == nil {
			log.V(2).Info("no rule")
			return
		}

		rule.Callback(w, r)
	})

	return mux
}
