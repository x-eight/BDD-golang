package main

import (
	"log"
	"time"
)

func run() error {
	var addr string = ":3000"
	server := NewServer(addr)

	client, err := server.SendAndClient()

	if err != nil {
		return err
	}

	//doGreet(client)
	//doGreetManyTimes(client)
	//doLongGreet(client)
	//doGreetEveryone(client)
	//doGreetWithDeadline(client, 5*time.Second)
	doGreetWithDeadline(client, 10*time.Second)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
