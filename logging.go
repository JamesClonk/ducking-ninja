package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
)

var (
	logFilename = "ducking-ninja.log"
)

func init() {
	log.SetOutput(w)
}

func logger() negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}
		log.Printf(";Request;%v;%v;%v", req.Method, req.URL.Path, addr)
		next(w, req)
	}
}
