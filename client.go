package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	CMDAddr = "localhost:4000"
)

type Client struct {
	nick    string
	conn    net.Conn
	history []string
}

func (client *Client) showChat() {
	fmt.Print("\033[2J")
	var buf bytes.Buffer
	buf.WriteString("##################################################################\n")
	for _, msg := range client.history {
		buf.WriteString(msg)
	}
	buf.WriteString("You: ")
	fmt.Print(buf.String())
}

func (client *Client) printIncoming() {
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Disconnected from server.", err)
			return
		} else if err != nil {
			fmt.Println("Could not get the message. Err:", err)
			return
		}
		client.history = append(client.history, message)
		client.showChat()
	}
}

func (client *Client) sendOutgoing(reader *bufio.Reader) {
	for {
		fmt.Print("You: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Bye!")
			return
		}

		_, err = fmt.Fprintf(client.conn, text)
		if err != nil {
			fmt.Println("Could not send the message. Err:", err)
			return
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", CMDAddr)
	if err != nil {
		fmt.Println("Could not connect to server. Err:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Enter nickname: ")
	reader := bufio.NewReader(os.Stdin)
	nick, err := reader.ReadString('\n')
	if err != nil {
		panic("lol")
	}

	nick = strings.TrimSpace(nick)
	_, err = fmt.Fprintf(conn, nick)
	if err != nil {
		fmt.Println("Could not send nick. Err:", err)
		panic("lol")
	}

	client := Client{nick, conn, make([]string, 0, 10)}

	fmt.Println("Your nickname is:", client.nick)

	go client.sendOutgoing(reader)
	client.printIncoming()
}
