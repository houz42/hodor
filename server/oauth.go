package server

import (
	"context"
	"net/http"

	"github.com/qiniu-ava/pkg/random"
	"golang.org/x/oauth2"
)

const (
	stateLen = 16
	stateTTL = 300 //seconds

	stateCookieKey = "hodor_state"
	tokenCookieKey = "hodor_token"
	originQueryKey = "hodor_origin"
)

type OAuthRule struct {
	Config   *oauth2.Config
	Upstream string
}

func (or *OAuthRule) Verify(w http.ResponseWriter, r *http.Request) {
	ck, e := r.Cookie(tokenCookieKey)
	if e != nil {
		or.redirectToLogin(w, r)
		return
	}
	token, err := unmarshalToken(ck.Value)
	if err != nil {
		or.redirectToLogin(w, r)
		return
	}
	newToken, err := or.Config.TokenSource(context.TODO(), token).Token()
	if err != nil {
		or.redirectToLogin(w, r)
		return
	}
	if token.AccessToken != newToken.AccessToken {
		or.setToken(w, r, newToken)
		return
	}

	// valid request shall pass
	w.WriteHeader(http.StatusOK)
}

func (or *OAuthRule) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	state := random.SecureRandomGenerator.MustGenString(stateLen, 62)
	redirection := or.Config.AuthCodeURL(state, oauth2.SetAuthURLParam(originQueryKey, r.URL.String()))

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieKey,
		Value:    state,
		Domain:   r.Host,
		MaxAge:   stateTTL,
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, redirection, http.StatusFound)
}

// set token into cookie, and retry the original request
func (or *OAuthRule) setToken(w http.ResponseWriter, r *http.Request, t *oauth2.Token) {
	tk, e := marshalToken(t)
	if e != nil {
		internalError(w, "write token failed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieKey,
		Value:    tk,
		Domain:   r.Host,
		Expires:  t.Expiry,
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, r.URL.String(), http.StatusFound)
}

func (or *OAuthRule) Callback(w http.ResponseWriter, r *http.Request) {
	queryState := r.URL.Query().Get("state")
	ck, e := r.Cookie(stateCookieKey)
	if e != nil {
		unauthorizedError(w, "no state in cookie")
	}
	if queryState != ck.Value {
		unauthorizedError(w, "state mismatch")
		return
	}

	token, err := or.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		unauthorizedError(w, "got no token")
		return
	}

	or.setToken(w, r, token)
}
