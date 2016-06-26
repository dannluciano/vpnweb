package main

import (
	"html/template"
	"net/http"
	"runtime"
	"os"
)

var templates = template.Must(template.ParseFiles("status.html", "env.html"))

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	BasicAuth(w, r, func() {
		clientList := parseStatusFile()
		err := templates.ExecuteTemplate(w, "status.html", &clientList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	if env == "dev" {
		templates = template.Must(template.ParseFiles("status.html", "env.html"))
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	http.Redirect(w, r, "/status/", http.StatusFound)
}

func EnvHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	BasicAuth(w, r, func() {

		envs := Envs{
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			os.Getenv("GOPATH"),
			runtime.GOROOT(),
			runtime.Compiler,
		}

		err := templates.ExecuteTemplate(w, "env.html", &envs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
