package main

import (
	"net"
	"log"
	"fmt"
	"io"
)

var (
	counter int
	listenAddr = "localhost:8080"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

// this code is from: https://www.youtube.com/watch?v=QTBZxDgRZM0&ab_channel=ahmetalpbalkan

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("failed to listen")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err)
		}

		backend := chooseBackend()
		fmt.Printf("counter: %d backend: %s\n", counter, backend)
		go func() {
			err := proxy(conn, backend)
			if err != nil {
				log.Printf("WARNING: proxying failed: %v", err)
			}
		} ()
	}
}

func proxy(c net.Conn, backend string) error {
	bc, err := net.Dial("tcp", backend)
	if err != nil {
		return fmt.Errorf("failed to connect to backend %s : %v", backend, err)
	}

	// c -> bc
	go io.Copy(bc, c)

	// bc -> c
	go io.Copy(c, bc)

	return nil
}

func chooseBackend() string {
	s := server[counter % len(server)]
	counter++
	return s 
}