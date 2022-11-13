package main

import (
	"fmt"
	"io"
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

	for {
		// Waiting for the connection
		conn, err := l.Accept()
		fmt.Println("Got connected", conn.RemoteAddr())
		if err != nil {
			fmt.Printf("Error accepting conns: %v", err)
			os.Exit(1)
		}

		go func(conn net.Conn) {
			defer conn.Close()
			msg := make([]byte, 4028)
			if _, err := conn.Read(msg); err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Println("Error reading from client: ", err.Error())
					os.Exit(1)
				}

			}
			fmt.Println("Got command: ", string(msg))
			conn.Write([]byte("+PONG\r\n"))
			fmt.Println("closing connection: ", conn.RemoteAddr())
		}(conn)
	}
}
