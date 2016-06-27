package main

import (
	"bufio"
	"log"
	"os"
	"strings"
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

type Client struct {
	CommonName     string `json:"Common Name"`
	RealAddress    string `json:"Real Address"`
	VirtualAddress string `json:"Virtual Address"`
	ConnectedSince string `json:"Connected Since"`
	LastRef        string `json:"Last Ref"`
}

type ClientList struct {
	Clients []Client `json:"Clients"`
	Status  string   `json:"Status"`
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

		client := Client{CommonName: row[0], RealAddress: row[1], ConnectedSince: row[4]}

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
