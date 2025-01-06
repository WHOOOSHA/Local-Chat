package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	userList map[string]net.Conn
	mutex    sync.Mutex
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) getName(conn net.Conn) (string, error) {
	var nameStr string
	for {
		conn.Write([]byte("Enter your name: "))
		nameReader := bufio.NewReader(conn)
		name, err := nameReader.ReadString('\n')
		if err != nil {
			return "", err
		}
		nameStr = name[:len(name)-1]

		s.mutex.Lock()
		if _, exists := s.userList[nameStr]; exists {
			conn.Write([]byte("This name already exists.\n"))
		} else {
			s.userList[nameStr] = conn
			conn.Write([]byte("Success! You may send messages!\n"))
			s.mutex.Unlock()
			break
		}
		s.mutex.Unlock()
	}
	return nameStr, nil
}

func (*Server) getMessage(conn net.Conn) (string, error) {
	messageReader := bufio.NewReader(conn)
	message, err := messageReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return message[:len(message)-1], nil
}

func (*Server) parseMessage(message string) (string, string) {
	if message[0] == '@' {
		for ind, ch := range message {
			if ch == ' ' {
				return message[1:ind], message[ind+1:]
			}
		}
	}
	return "", message
}

func (s *Server) sendMessage(name string, recipientName string, message string) {
	if recipientName != "" {
		recipientConn := s.userList[recipientName]
		_, err := recipientConn.Write([]byte("New message from " + name + ": " + message + "\n"))
		if err != nil {
			s.mutex.Lock()
			delete(s.userList, recipientName)
			s.mutex.Unlock()
		}
	} else {
		for n, recipientConn := range s.userList {
			if n == name {
				continue
			}
			_, err := recipientConn.Write([]byte("New message from " + name + ": " + message + "\n"))
			if err != nil {
				s.mutex.Lock()
				delete(s.userList, n)
				s.mutex.Unlock()
			}
		}
	}
}

func (s *Server) handlerConn(conn net.Conn) {
	defer conn.Close()
	name, err := s.getName(conn)
	if err != nil {
		return
	}

	for {
		message, err := s.getMessage(conn)
		if err != nil {
			s.mutex.Lock()
			delete(s.userList, name)
			s.mutex.Unlock()
			return
		}

		recipientName, text := s.parseMessage(message)
		s.sendMessage(name, recipientName, text)
	}
}

func (s *Server) run() {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server running on", addr)
	s.userList = make(map[string]net.Conn)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handlerConn(conn)
	}
}
