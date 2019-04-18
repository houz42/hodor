package server

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthRule struct {
	Config   *oauth2.Config
	Upstream string
}

//
func (or *OAuthRule) Verify(w http.ResponseWriter, r *http.Request) {
	token, err := or.retriveToken(r)
	if err != nil {
		failed(w)
		return
	}
	newToken, err := or.Config.TokenSource(context.TODO(), token).Token()
	if err != nil {
		failed(w)
		return
	}
	if token.AccessToken != newToken.AccessToken {
		or.setToken(w, r, newToken)
	}

}

// redirect to login
func (or *OAuthRule) Login(w http.ResponseWriter, r *http.Request) {
	// todo: state will be set into cookie
}

func (or *OAuthRule) retriveToken(r *http.Request) (*oauth2.Token, error) {
	// todo: retrive from cookie
	return nil, nil
}

func (or *OAuthRule) setToken(w http.ResponseWriter, r *http.Request, t *oauth2.Token) {
	// todo: set token into cookie, and redirect to retry
}

func (or *OAuthRule) Callback(w http.ResponseWriter, r *http.Request) {
	if err := checkState(r); err != nil {
		failed(w)
		return
	}

	token, err := or.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		failed(w)
		return
	}

	token.SetAuthHeader(r)
	w.WriteHeader(http.StatusOK)
}

// compare state set in cookie with the one in callback
func checkState(r *http.Request) error {
	return nil
}

func failed(w http.ResponseWriter) {}
