package main

import (
	"log"

	"github.com/frzifus/nuki"
)

var version string

func main() {
	log.Printf("Version %s\n", version)
	n := nuki.NewNuki("127.0.0.1", "8080")
	if _, err := n.Auth(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Authentication successful, your token is:", n.Token())
}
