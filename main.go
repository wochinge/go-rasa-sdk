package main

import (
	"fmt"
	"github.com/wochinge/go-rasa-sdk/constants"
	"github.com/wochinge/go-rasa-sdk/server"
	"log"
)

func main() {
	fmt.Println("Runing Rasa action server 🚀")
	fmt.Printf("%s", "dasd")
	log.Fatal(server.Serve(constants.DefaultServerPort))
	fmt.Println("Goodbye ...")
}
