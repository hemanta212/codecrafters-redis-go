package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("Error binding tcp listener: %v", err)
		os.Exit(1)
	}
	fmt.Println(":: Listening on port 6379.. ")
	_, err = l.Accept()
	if err != nil {
		fmt.Printf("Error accepting conns: %v", err)
		os.Exit(1)
	}
}
