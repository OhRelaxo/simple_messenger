package chatmanager

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type setName struct {
	conn net.Conn
	name string
}

type getName struct {
	conn net.Conn
	resp chan string
}

type listNames struct {
	resp chan []string
}

type ChatManager struct {
	manager chan any
}

func (cm *ChatManager) setName(conn net.Conn, name string) {
	cm.manager <- setName{
		conn: conn,
		name: name,
	}
}

func (cm *ChatManager) getName(conn net.Conn) string {
	resp := make(chan string)
	cm.manager <- getName{
		conn: conn,
		resp: resp,
	}
	return <-resp
}

func (cm *ChatManager) listNames() []string {
	resp := make(chan []string)
	cm.manager <- listNames{
		resp: resp,
	}
	return <-resp
}

func (cm *ChatManager) run() {
	clients := make(map[net.Conn]string)

	for msg := range cm.manager {
		switch m := msg.(type) {
		case setName:
			clients[m.conn] = m.name
		case getName:
			m.resp <- clients[m.conn]
		case listNames:
			names := make([]string, 0, len(clients))
			for _, name := range clients {
				names = append(names, name)
			}
			m.resp <- names
		}
	}
}

func NewChatManager() *ChatManager {
	m := ChatManager{
		manager: make(chan any),
	}
	go m.run()
	return &m
}

func (cm *ChatManager) HandleConnection(conn net.Conn) {
	defer conn.Close()

	err := sendResponse("write your name:\n", conn)
	if err != nil {
		log.Println("failed to send response: ", err)
	}
	name, err := reciveMessage(conn)
	if err != nil {
		log.Println("failed to read message: ", err)
		return
	}
	cm.setName(conn, name)
	nameManager := cm.getName(conn)
	if nameManager != name {
		log.Printf("ein Client hat ein falschen Namen erhalten\nname laut Manager: %s\nname laut Client: %s\n", nameManager, name)
	}
	err = sendResponse("dein Name ist: "+name, conn)
	if err != nil {
		log.Println("failed to send response: ", err)
		return
	}

	for {
		// hier soll es dann eine Möglichkeit geben Broadcast, multicast und unicast zu nutzen (routing), sowie sich die aktuellen nutzer anzeigen zu lassen,
		// um zu entscheiden mit welchen benutzer man sich unterhalten möchte.

		message, err := reciveMessage(conn)
		if err != nil {
			log.Println("failed to read message: ", err)
			break
		}
		fmt.Print("Received message: ", message)

		cleaned := cleanUpMessage(message)
		log.Printf("cleaned Message: '%s'", cleaned)

		var response string
		switch cleaned {
		case "help":
			response += "du kannst die befehle:\nlistclients\nausführen\n"
		case "listclients":
			clients := cm.listNames()
			for _, client := range clients {
				if client == nameManager {
					trimed := strings.TrimRight(client, "\n")
					response += trimed + " (you)" + "\n"
					continue
				}
				response += client + "\n"
			}
		default:
			response += "Received message: " + message
		}

		log.Println("responding to client with response: ", response)
		err = sendResponse(response, conn)
		if err != nil {
			log.Println("failed to send response: ", err)
			return
		}
	}
}

func cleanUpMessage(message string) string {
	trimOne := strings.Trim(message, " ")
	trimeTwo := strings.TrimRight(trimOne, "\n")

	lower := strings.ToLower(trimeTwo)
	return lower
}

func sendResponse(message string, conn net.Conn) error {
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

func reciveMessage(conn net.Conn) (string, error) {
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

/*
 * beispiel code für Broadcast
 *
func broadcast(sender net.Conn, msg string) {
     for conn := range clients {
         if conn != sender {
             fmt.Fprintln(conn, msg)
        }
    }
}

*/
