package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (*Client) handleServerMessages(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection with server lost.")
			os.Exit(1)
		}
		fmt.Println(strings.TrimSpace(message))
	}
}

func (*Client) handleUserInput(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Failed to send message:", err)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func (c *Client) run() {
	addr := ":8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic("Cannot connect to server")
	}
	defer conn.Close()

	message := make([]byte, 1024)
	_, err = conn.Read(message)
	if err != nil {
		panic("Error reading from server")
	}
	fmt.Println(strings.TrimSpace(string(message)))
	go c.handleServerMessages(conn)

	c.handleUserInput(conn)
}
