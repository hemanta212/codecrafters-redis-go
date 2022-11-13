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
	defer l.Close()
	fmt.Println(":: Listening on port 6379.. ")
	// Waiting for the connection
	for {
		conn, err := l.Accept()
		fmt.Println("Got connected", conn.RemoteAddr())
		if err != nil {
			fmt.Printf("Error accepting conns: %v", err)
			os.Exit(1)
		}
		msg := make([]byte, 4028)
		len, _ := conn.Read(msg)
		fmt.Println("Got command: ", string(msg[:len]))
		conn.Write([]byte("+PONG\r\n"))
		fmt.Println("closing connection: ", conn.RemoteAddr())
		conn.Close()
	}
}
