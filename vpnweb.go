package main

import (
	"os"
	"log"
	"net/http"
)

var (
	statusFile = "/etc/openvpn/openvpn-status.log"
	username = "dannluciano"
	password = "dlcorp"
	env = "dev"
)

type Envs struct {
	GoVersion string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GOPATH    string
	GOROOT    string
	Compiler  string
}

func SetupEnv() {
	environment := os.Getenv("ENV")
	if environment != "" {
		env = environment
	}
	log.Println("ENV=", env)
	if env == "dev" {
		statusFile = "openvpn-status.log"
	}
}

func main() {
	SetupEnv()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/env/", EnvHandler)
	http.HandleFunc("/reload", ReloadHandler)
	http.HandleFunc("/status/", StatusHandler)

	http.Handle("/static/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", nil)
}