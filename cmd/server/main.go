package main

import (
	"log"
	"net"

	chatmanager "github.com/OhRelaxo/tcp_test/cmd/server/internal/chatManager"
)

const (
	ip   = "localhost"
	port = "8080"
)

func main() {
	log.Println("creating manager routine...")
	manager := chatmanager.NewChatManager()
	log.Println("successfully created manager routine")

	log.Printf("listening on ip: %v and port: %v\n", ip, port)
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		log.Fatalf("fatal error while listining on IP: %v and Port: %v : %v", ip, port, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to Accept message: %v\n", err)
		}

		go manager.HandleConnection(conn)
	}
}
