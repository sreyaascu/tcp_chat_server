package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

var clients []net.Conn
var mu sync.Mutex

func removeClient(owner string, conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for i, client := range clients {
		if conn == client {
			fmt.Printf("%s left\n", owner)
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func handleClient(owner string, conn net.Conn) {
	defer conn.Close()
	var prefix string
	buff := make([]byte, 2048)
	for {
		prefix = owner + ": "
		n, read_err := conn.Read(buff)
		if read_err != nil {
			fmt.Println("Reading Error")
			removeClient(owner, conn)
			return
		}
		msg := string(buff[:n])
		prefix += msg
		if strings.TrimSpace(msg) == "stop" {
			removeClient(owner, conn)
			return
		}
		mu.Lock()
		local_clients := make([]net.Conn, len(clients))
		copy(local_clients, clients)
		mu.Unlock()
		for _, client := range local_clients {
			if conn == client {
				continue
			}
			client.Write([]byte(prefix))
		}
		prefix = ""
		fmt.Printf("%s: %s", owner, msg)

	}
}

func main() {
	id := 0
	listener, list_err := net.Listen("tcp", ":8080")
	if list_err != nil {
		fmt.Println("Listen Error: ", list_err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port 8080")
	fmt.Println("Waiting for connection...")
	for {
		conn, acc_err := listener.Accept()
		if acc_err != nil {
			fmt.Println("Accept Error: ", acc_err)
			return
		}
		mu.Lock()
		clients = append(clients, conn)
		mu.Unlock()
		// defer conn.Close()
		fmt.Printf("term %d joined\n", id)
		go handleClient("term "+strconv.Itoa(id), conn)
		id++

	}
}
