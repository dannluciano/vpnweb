package main

import (
	"log"
	"net/http"
)

func Logger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " " + r.URL.String())
}

func BasicAuth(w http.ResponseWriter, r *http.Request, lambda func()) {
	if r.Header.Get("Authorization") == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="mydomain"`)
	} else {
		u, p, ok := r.BasicAuth()
		if ok && username == u && password == p {
			lambda()
			return
		}
	}
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}
