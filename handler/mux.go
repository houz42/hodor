package handler

import (
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
)

type Config struct {
	Rules        map[string]*OAuthRule `json:"rueles,omitempty"`
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
		origin := r.URL.Query().Get(originQueryKey)
		if origin == "" {
			log.V(3).Info("no origin url in callback")
			badRequest(w, "no origin url in callback")
			return
		}
		ou, e := url.Parse(origin)
		if e != nil {
			log.Error(e, "parse origin url in callback request failed")
			badRequest(w, "bad origin url in callback")
			return
		}
		rule := cfg.Rules[ou.Host]
		if rule == nil {
			log.V(2).Info("no rule")
			return
		}

		rule.Callback(w, r)
	})

	return mux
}
