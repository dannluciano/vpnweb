package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"strings"
	"net/http"
	"html/template"
	"runtime"
)

var (
	statusFile = "/etc/openvpn/openvpn-status.log"
	username = "dannluciano"
	password = "dlcorp"
	env = "dev"
)

type Client struct {
	CommonName     string `json:"Common Name"`
	RealAddress    string `json:"Real Address"`
	VirtualAddress string `json:"Virtual Address"`
	ConnectedSince string `json:"Connected Since"`
	LastRef        string `json:"Last Ref"`
}

func (c *Client) String() string {
	return fmt.Sprintf(
		"Common Name: %s, Real Address: %s, Virtual Address: %s, Connected Since: %s, Last Ref: %s",
		c.CommonName, c.RealAddress, c.VirtualAddress, c.ConnectedSince, c.LastRef)
}

type ClientList struct {
	Clients []Client        `json:"Clients"`
	Status  string                `json:"Status"`
}

type Envs struct {
	GoVersion string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GOPATH    string
	GOROOT    string
	Compiler  string
}

func parseStatusFile() ClientList {

	file, err := os.Open(statusFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	_, err = reader.ReadString('\n')
	if err != err {
		log.Fatal(err)
	}

	status, err := reader.ReadString('\n')
	if err != err {
		log.Fatal(err)
	}
	status = strings.TrimSpace(status)

	_, err = reader.ReadString('\n')
	if err != err {
		log.Fatal(err)
	}

	clients := make([]Client, 0)

	for true {

		line, err := reader.ReadString('\n')
		if err != err {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)

		if strings.Contains(line, "ROUTING TABLE") {
			break
		}

		row := strings.Split(line, ",")
		if len(row) != 5 {
			log.Fatal("Parse Failed. Expected 5 Fields in CLIENT LIST")
		}

		client := Client{CommonName:row[0], RealAddress:row[1], ConnectedSince:row[4]}

		clients = append(clients, client)
	}

	_, err = reader.ReadString('\n')
	if err != err {
		log.Fatal(err)
	}

	for _, _ = range clients {
		line, err := reader.ReadString('\n')
		if err != err {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)

		row := strings.Split(line, ",")
		if len(row) != 4 {
			log.Fatal("Parse Failed. Expected 4 Fields in ROUTING TABLE goted ", len(row))
		}

		for index, client := range clients {
			if strings.Compare(client.CommonName, row[1]) == 0 {
				client.VirtualAddress = row[0]
				client.LastRef = row[3]

				clients[index] = client
			}
		}
	}

	return ClientList{clients, status}
}

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