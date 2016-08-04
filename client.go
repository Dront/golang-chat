package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	CMDAddr = "localhost:4000"
)

func printIncoming(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Disconnected from server.", err)
			return
		} else if err != nil {
			fmt.Println("Could not get the message. Err:", err)
			return
		}
		fmt.Print(message)
	}
}

func sendOutgoing(conn net.Conn, reader *bufio.Reader) {
	for {
		fmt.Print("You: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Bye!")
			break
		}

		_, err = fmt.Fprintf(conn, text)
		if err != nil {
			fmt.Println("Could not send the message. Err:", err)
			break
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

	fmt.Println("Your nickname is:", nick)

	go sendOutgoing(conn, reader)
	printIncoming(conn)
}
