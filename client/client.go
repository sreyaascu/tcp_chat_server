package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func receiveMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println(msg)

	}
}

func readKeyBoard(conn net.Conn) {
	var msg string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg = scanner.Text()
		conn.Write([]byte(msg + "\n"))
		if msg == "stop" {
			return
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Dial error: ", err)
		return
	}
	defer conn.Close()
	go receiveMessage(conn)

	readKeyBoard(conn)

}
