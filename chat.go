package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Message struct {
	nick string
	data string
}

func (msg Message) String() string {
	return msg.nick + " : " + msg.data
}

type Client struct {
	nick     string
	conn     net.Conn
	toClient chan Message
}

func (c Client) ReadLinesInto(msgChan chan Message) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		msgChan <- Message{c.nick, strings.TrimSpace(line)}
	}
}

func (c Client) SendIncoming() {
	for msg := range c.toClient {
		if msg.nick != c.nick {
			_, err := io.WriteString(c.conn, fmt.Sprintf("%s: %s\n", msg.nick, msg.data))
			if err != nil {
				return
			}
		}
	}
}

func handleNewClient(c net.Conn, addClientChan, delClientChan chan Client, msgChan chan Message) {
	addr := c.RemoteAddr().String()
	logger := log.New(os.Stdin, addr+" ", log.Ltime)
	buf := make([]byte, 4096)
	logger.Printf("Connection opened.")
	defer func() {
		c.Close()
		logger.Printf("Connection closed")
	}()

	n, err := c.Read(buf)
	if err != nil || n == 0 {
		logger.Printf("Could not get nickname. Err:", err)
		panic("lol")
	}

	client := Client{string(buf[0:n]), c, make(chan Message)}

	addClientChan <- client
	defer func() {
		delClientChan <- client
	}()

	go client.SendIncoming()
	client.ReadLinesInto(msgChan)
}

func chatRoom(addClientChan, delClientChan chan Client, msgChan chan Message) {
	logger := log.New(os.Stdin, "chat room ", log.Ltime)
	clients := make(map[string]Client)

	for {
		select {
		case msg := <-msgChan:
			logger.Printf("%s: %s\n", msg.nick, msg.data)
			for _, client := range clients {
				client.toClient <- msg
			}
		case newClient := <-addClientChan:
			logger.Println("new client:", newClient.nick)
			for _, client := range clients {
				client.toClient <- Message{newClient.nick, "connected"}
			}
			clients[newClient.nick] = newClient
		case delClient := <-delClientChan:
			delete(clients, delClient.nick)
			for _, client := range clients {
				client.toClient <- Message{delClient.nick, "disconnected"}
			}
			logger.Println("removed client:", delClient.nick)
		}
	}
}

func main() {
	logger := log.New(os.Stdin, "Main ", log.Ltime)
	const CMDAddr = "localhost:4000"
	logger.Println("Launching server...")

	listener, err := net.Listen("tcp", CMDAddr)
	if err != nil {
		logger.Fatal(err)
	}
	defer listener.Close()

	msgChan := make(chan Message)
	addClientChan := make(chan Client)
	delClientChan := make(chan Client)
	go chatRoom(addClientChan, delClientChan, msgChan)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println(err)
		}
		go handleNewClient(conn, addClientChan, delClientChan, msgChan)
	}
}
