package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
)

type Auth struct {
	Username string
	Password string
}

func readAuth(filename string) (auth Auth) {
	readJson(filename, &auth)
	return
}

func compareSHA(a string, b string) bool {
	sha1 := sha256.Sum256([]byte(a))
	sha2 := sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(sha1[:], sha2[:]) == 1
}

func authenticate(r *render.Render, a Auth) negroni.HandlerFunc {
	var secret = base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password))

	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		auth := req.Header.Get("Authorization")
		if !compareSHA("Basic "+secret, auth) {
			w.Header().Set("WWW-Authenticate", `Basic realm="ducking-ninja does not like you yet!"`)
			r.JSON(w, http.StatusUnauthorized, "I don't like you!")
			return
		}
		w.Header().Set("X-Authenticated-Username", a.Username)
		next(w, req)
	}
}
