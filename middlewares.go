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
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	} else {
		u, p, ok := r.BasicAuth()
		log.Println(username, password)
		if ok {
			if username == u && password == p {
				lambda()
			} else {
				http.Error(w, "Not Authorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
		}
	}
}