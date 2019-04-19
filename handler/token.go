package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func marshalToken(t *oauth2.Token) (string, error) {
	buf := &bytes.Buffer{}
	if e := gob.NewEncoder(buf).Encode(t); e != nil {
		return "", errors.Wrap(e, "encode token failed")
	}
	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	return s, nil
}

func unmarshalToken(s string) (*oauth2.Token, error) {
	buf, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return nil, errors.Wrap(e, "parse token string failed")
	}
	t := &oauth2.Token{}
	if e := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(t); e != nil {
		return nil, errors.Wrap(e, "decode token failed")
	}
	return t, nil
}
