package chatmanager

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	reader := bufio.NewReader(conn)

	_, err := conn.Write([]byte("write your name:\n"))
	if err != nil {
		log.Println("failed to send response: ", err)
	}
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println("failed to read message: ", err)
		return
	}
	cm.setName(conn, name)
	nameManager := cm.getName(conn)
	if nameManager != name {
		log.Printf("ein Client hat ein falschen Namen erhalten\nname laut Manager: %s\nname laut Client: %s\n", nameManager, name)
	}
	_, err = conn.Write([]byte("dein Name ist: " + name))
	if err != nil {
		log.Println("failed to send response: ", err)
		return
	}

	/*
		allClients := cm.listNames()
		for _, client := range allClients {

		}
	*/

	for {
		// hier soll es dann eine Möglichkeit geben Broadcast, multicast und unicast zu nutzen (routing), sowie sich die aktuellen nutzer anzeigen zu lassen,
		// um zu entscheiden mit welchen benutzer man sich unterhalten möchte

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read message: ", err)
			break
		}
		fmt.Print("Received message: ", message)

		_, err = conn.Write([]byte("got message: " + message))
		if err != nil {
			log.Println("failed to send response: ", err)
			break
		}
	}
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
