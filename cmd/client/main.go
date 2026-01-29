package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	ip   = "localhost"
	port = "8080"
)

func main() {
	// Use net.JoinHostPort for IPv6 safety.
	address := net.JoinHostPort(ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("fatal error while trying to dial up: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server, checking for message from server...")
	response, _ := bufio.NewReader(conn).ReadString('\n')
	if response != "" {
		fmt.Printf("Response from server: \n%s", response)
	} else {
		fmt.Println("no message from server, send a message:")
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		message, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Error receiving response:", err)
			break
		}
		fmt.Printf("Response from server: \n%s", response)
	}
}
