package main

import (
	"html/template"
	"net/http"
	"os"
	"runtime"
	"log"
)

var (
	statusTmpl *template.Template
	envTmpl *template.Template
)

func LoadTemplates() {
	log.Println("Load Templates")

	statusTmpl = template.Must(template.New("status").ParseFiles(
		"templates/base.html",
		"templates/status.html"))
	envTmpl = template.Must(template.New("env").ParseFiles(
		"templates/base.html",
		"templates/env.html"))

}

func RenderTemplate(w http.ResponseWriter, tmpl*template.Template, p interface{}) {
	err := tmpl.ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	BasicAuth(w, r, func() {
		clientList := parseStatusFile()
		RenderTemplate(w, statusTmpl, &clientList)
	})
}

func EnvHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	BasicAuth(w, r, func() {

		envs := Envs{
			os.Getenv("ENV"),
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			os.Getenv("GOPATH"),
			runtime.GOROOT(),
			runtime.Compiler,
		}

		RenderTemplate(w, envTmpl, &envs)
	})
}

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	Setup()

	http.Redirect(w, r, "/", http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Logger(w, r)

	http.Redirect(w, r, "/status/", http.StatusFound)
}
