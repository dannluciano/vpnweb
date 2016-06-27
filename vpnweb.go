package main

import (
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"os"
)

var (
	statusFile = "/etc/openvpn/openvpn-status.log"
	username   = "dannluciano"
	password   = "dlcorp"
	env        = "dev"
)

func SetupEnv() {
	environment := os.Getenv("ENV")
	if environment != "" {
		env = environment
	}
	log.Println("ENV=", env)
	if env == "dev" {
		statusFile = "public/openvpn-status.log"
	}
}

func main() {
	SetupEnv()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/env/", EnvHandler)
	http.HandleFunc("/reload", ReloadHandler)
	http.HandleFunc("/status/", StatusHandler)

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(cwd+"/public"))))

	srv := &http.Server{
		Addr: ":8000",
	}
	http2.ConfigureServer(srv, &http2.Server{})
	log.Fatal(srv.ListenAndServeTLS(cwd+"/keys/cert.pem", cwd+"/keys/key.pem"))
}
