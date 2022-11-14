package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
	storage := NewKeyValueStore()

	for {
		// Waiting for the connection
		conn, err := l.Accept()
		fmt.Println(":: Got connected", conn.RemoteAddr())
		if err != nil {
			fmt.Printf(":: Error accepting conns: %v", err)
			os.Exit(1)
		}

		go handleConn(conn, storage)
	}
}

func handleConn(conn net.Conn, storage *KeyValueStore) {
	defer conn.Close()
	for {
		value, err := decodeRESP(conn)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Println(":: Error reading from client: ", err.Error())
				return
			}
		}
		responseHandler(value, conn, storage)
	}
}

func decodeRESP(conn io.Reader) ([]string, error) {
	msg := make([]byte, 1024)
	value := []string{}
	msglen, err := conn.Read(msg)
	if err != nil {
		return value, err
	}
	message := strings.TrimSpace(string(msg[:msglen]))
	messageSlice := strings.Split(message, "\r\n")
	fmt.Println(":: Message content: ", messageSlice)
	// discard the first two RESP spec keywords and pick out the commands at even interval
	for i, item := range messageSlice[2:] {
		if i%2 == 0 {
			value = append(value, item)
		}
	}
	fmt.Println(":: Parsed: ", value)
	return value, nil
}

func responseHandler(instructions []string, conn io.Writer, storage *KeyValueStore) {
	command, args := strings.ToLower(instructions[0]), instructions[1:]
	switch command {
	case "echo":
		echoedWord := args[0]
		output := fmt.Sprintf("$%d\r\n%s\r\n", len(echoedWord), echoedWord)
		// fmt.Println(":: Writing result as: ", strconv.Quote(output))
		conn.Write([]byte(output))
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	case "set":
		setHandler(args, conn, storage)
	case "get":
		getHandler(args, conn, storage)
	default:
		conn.Write([]byte("+PONG\r\n"))
	}
}

func setHandler(args []string, conn io.Writer, storage *KeyValueStore) {
	key, value := args[0], args[1]
	if len(args) > 3 {
		px := args[3]
		duration, _ := strconv.Atoi(px)
		go time.AfterFunc(time.Duration(duration)*time.Millisecond, func() { storage.Delete(key) })
	}
	storage.Set(key, value)
	conn.Write([]byte("+OK\r\n"))
}

func getHandler(args []string, conn io.Writer, storage *KeyValueStore) {
	value, found := storage.Get(args[0])
	if found {
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
	} else {
		conn.Write([]byte("$-1\r\n"))
	}
}
