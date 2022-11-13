package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
		fmt.Println(":: Got connected", conn.RemoteAddr())
		if err != nil {
			fmt.Printf(":: Error accepting conns: %v", err)
			os.Exit(1)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		msg := make([]byte, 1024)
		msglen, err := conn.Read(msg)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Println(":: Error reading from client: ", err.Error())
				os.Exit(1)
			}

		}
		command := strings.TrimSpace(string(msg[:msglen]))
		fmt.Println(":: Got command: ", strconv.Quote(command))

		if strings.Contains(string(msg), "ECHO") {
			stringSlice := strings.Split(command, "\r\n")
			echoedWord := stringSlice[4]
			// fmt.Println(":: List: ", stringSlice)

			output := fmt.Sprintf("+%s\r\n", echoedWord)
			// fmt.Println(":: Writing result as: ", strconv.Quote(output))
			conn.Write([]byte(output))
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
