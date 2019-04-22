package handler

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/qiniu-ava/pkg/random"
	"golang.org/x/oauth2"
)

const (
	stateLen = 16
	stateTTL = 300 //seconds

	stateCookieKey = "hodor_state"
	tokenCookieKey = "hodor_token"
	hostCookieKey  = "hodor_host"
	pathCookieKey  = "hodor_path"
)

type OAuthRule struct {
	Config *oauth2.Config
	log    logr.Logger
}

func (or *OAuthRule) Verify(w http.ResponseWriter, r *http.Request) {
	l := or.log.WithValues("host", r.Host, "path", r.URL.Path, "method", r.Method)
	ck, e := r.Cookie(tokenCookieKey)
	if e != nil {
		l.V(4).Info("request without cookie, redirect to login")
		or.redirectToLogin(w, r)
		return
	}
	token, err := unmarshalToken(ck.Value)
	if err != nil {
		l.V(4).Info("invalid token in cookie, redirect to login")
		or.redirectToLogin(w, r)
		return
	}

	// fixme: it's not safe trusting token in request cookie
	// the token should be checked over id provider or signed
	newToken, err := or.Config.TokenSource(context.TODO(), token).Token()
	if err != nil {
		l.Error(err, "token rotation failed, redirect to login")
		or.redirectToLogin(w, r)
		return
	}
	if token.AccessToken != newToken.AccessToken {
		l.V(4).Info("token refreshed, reload")
		or.setToken(w, r, newToken)
		return
	}

	// valid request shall pass
	l.V(3).Info("valid request")
	w.WriteHeader(http.StatusOK)
}

func (or *OAuthRule) Callback(w http.ResponseWriter, r *http.Request) {
	l := or.log.WithValues("host", r.Host, "path", r.URL.Path, "method", r.Method)
	queryState := r.URL.Query().Get("state")
	ck, e := r.Cookie(stateCookieKey)
	if e != nil {
		l.V(3).Info("no state in callback, reject")
		unauthorized(w, "no state in cookie")
		return
	}
	if queryState != ck.Value {
		l.V(2).Info("state mismatch, reject")
		unauthorized(w, "state mismatch")
		return
	}

	token, err := or.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		l.Error(err, "exchange for")
		unauthorized(w, "got no token")
		return
	}

	or.setToken(w, r, token)
}

func (or *OAuthRule) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	state := random.SecureRandomGenerator.MustGenString(stateLen, 62)
	redirection := or.Config.AuthCodeURL(state)

	setTempCookie(w, r, pathCookieKey, r.URL.Path)
	setTempCookie(w, r, hostCookieKey, r.Host)
	setTempCookie(w, r, stateCookieKey, state)
	http.Redirect(w, r, redirection, http.StatusFound)
}

// set token into cookie, and reload the original request
func (or *OAuthRule) setToken(w http.ResponseWriter, r *http.Request, t *oauth2.Token) {
	tk, e := marshalToken(t)
	if e != nil {
		or.log.Error(e, "marshal token failed")
		internalError(w, "write token failed")
		return
	}

	path := r.URL.Path
	oc, e := r.Cookie(pathCookieKey)
	if e != nil {
		if e != http.ErrNoCookie {
			or.log.Error(e, "get origin path in cookie failed")
			internalError(w, "")
			return
		}
	} else if oc.Value != "" {
		path = oc.Value
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieKey,
		Value:    tk,
		Domain:   r.Host,
		Expires:  t.Expiry,
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, path, http.StatusFound)
}
