package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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
	response, err := reciveResponse(conn)
	if response != "" {
		fmt.Printf("Response from server: \n%s", response)
	} else {
		fmt.Println("no message from server, send a message:")
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		message, _ := reader.ReadString('\n')
		err := sendMessage(message, conn)
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}

		response, err := reciveResponse(conn)
		if err != nil {
			log.Println("Error receiving response:", err)
			break
		}
		fmt.Printf("Response from server: \n%s", response)
	}
}

func reciveResponse(conn net.Conn) (string, error) {
	var length int32
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

func sendMessage(message string, conn net.Conn) error {
	msgBytes := []byte(message)

	length := int32(len(msgBytes))
	err := binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		return err
	}

	_, err = conn.Write(msgBytes)
	if err != nil {
		return err
	}
	return nil
}

//
