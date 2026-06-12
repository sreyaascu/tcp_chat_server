package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mu sync.Mutex

func renameUser(name string, conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	clients[conn] = name
}

func listUsers(conn net.Conn) {
	mu.Lock()
	names := make([]string, 0, len(clients))
	for _, name := range clients {
		names = append(names, name)
	}
	mu.Unlock()
	conn.Write([]byte("Connected Users: \n"))
	for _, name := range names {
		conn.Write([]byte(name + "\n"))
	}
}

func showFunctions(conn net.Conn) {
	conn.Write([]byte("Chat Controls:\n"))
	conn.Write([]byte("/help - show chat commands\n"))
	conn.Write([]byte("/list - list members in the server\n"))
	conn.Write([]byte("/quit - disconnect from the server\n"))
}

func removeClient(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		if conn == client {
			fmt.Printf("%s left\n", clients[conn])
			delete(clients, conn)
			break
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	var prefix string
	var username string
	first_msg := true
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if first_msg {
			username = strings.TrimSpace(msg)
			mu.Lock()
			clients[conn] = username
			mu.Unlock()
			first_msg = false
			fmt.Printf("%s joined the chat\n", username)
			continue
		}
		if strings.TrimSpace(msg) == "/quit" {
			removeClient(conn)
			return
		}
		if strings.TrimSpace(msg) == "/list" {
			listUsers(conn)
			continue
		}
		if strings.TrimSpace(msg) == "/help" {
			showFunctions(conn)
			continue
		}
		parts := strings.Fields(msg)
		if len(parts) >= 2 && parts[0] == "/name" {
			newName := parts[1]
			renameUser(newName, conn)
			continue
		}
		prefix = username + ": " + msg
		mu.Lock()
		connections := make([]net.Conn, 0, len(clients))
		for conn := range clients {
			connections = append(connections, conn)
		}
		mu.Unlock()
		//Writing
		for _, client := range connections {
			if conn == client {
				continue
			}
			client.Write([]byte(prefix + "\n"))
		}
		fmt.Printf("%s: %s\n", username, msg)

	}
}

func main() {
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
		go handleClient(conn)

	}
}
