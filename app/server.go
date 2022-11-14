package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

		message := strings.TrimSpace(string(msg[:msglen]))
		value := decodeRESP(message)
		command, args := strings.ToLower(value[0]), value[1:]

		if command == "echo" {
			echoedWord := args[0]
			output := fmt.Sprintf("+%s\r\n", echoedWord)
			// fmt.Println(":: Writing result as: ", strconv.Quote(output))
			conn.Write([]byte(output))
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}

func decodeRESP(message string) []string {
	messageSlice := strings.Split(message, "\r\n")
	fmt.Println(":: Message content: ", messageSlice)
	value := []string{}
	// discard the first two RESP spec keywords and pick out the commands at even interval
	for i, item := range messageSlice[2:] {
		if i%2 == 0 {
			value = append(value, item)
		}
	}
	fmt.Println(":: Parsed: ", value)
	return value
}
