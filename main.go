package main

import "log"

//eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	server := NewServer(":3000")
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
