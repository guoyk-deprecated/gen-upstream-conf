package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

const socket_path = "/var/run/docker.sock"
const env_key = "DOCKER_HOST"

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		println("Loop entered")
		n, err := r.Read(buf[:])
		if err != nil {
			panic(err)
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	docker_host := os.Getenv(env_key)

	var network = "unix"
	var addr = socket_path

	if len(docker_host) == 0 {
		fmt.Printf("Docker ENV \"DOCKER_HOST\" not found, try default socket path: \"%s\".\n", socket_path)
	} else {
		fmt.Printf("Docker ENV \"DOCKER_HOST\" found: \"%s\".\n", docker_host)
		comp := strings.Split(docker_host, "://")
		network = comp[0]
		addr = comp[1]
	}

	var conn, err = net.Dial(network, addr)

	if err != nil {
		fmt.Printf("!! Cannot connect \"%s\".\n", addr)
		return
	}

	req, err := http.NewRequest("GET", "/containers/json", nil)

	if err != nil {
		panic(err)
	}

	err = req.Write(conn)

	if err != nil {
		panic(err)
	}

	go reader(conn)

	var input string
	fmt.Scanln(&input)
}
