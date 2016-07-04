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
	cwd	   = "."
)

func Setup() {
	log.SetPrefix("VPNWEB: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CWD:", cwd)

	environment := os.Getenv("ENV")
	if environment == "" {
		os.Setenv("ENV", "dev")
	} else {
		env = environment
	}
	log.Println("ENV:", env)

	if env == "dev" {
		statusFile = "public/openvpn-status.log"
	}

	LoadTemplates()
}

func main() {
	Setup()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/env/", EnvHandler)
	http.HandleFunc("/reload/", ReloadHandler)
	http.HandleFunc("/status/", StatusHandler)

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(cwd + "/public"))))

	srv := &http.Server{
		Addr: ":8000",
	}
	http2.ConfigureServer(srv, &http2.Server{})
	log.Println("https://localhost:8000")
	log.Fatal(srv.ListenAndServeTLS(cwd + "/keys/cert.pem", cwd + "/keys/key.pem"))
}
