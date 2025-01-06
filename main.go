package main

import (
	"os"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("Arguments error")
	}

	switch args[1] {
	case "IsServer":
		server := NewServer()
		server.run()
	case "IsClient":
		client := NewClient()
		client.run()
	}

}
