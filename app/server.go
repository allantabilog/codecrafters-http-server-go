package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
func prepareEchoResponse(message string) string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(message)) + "\r\n\r\n" + message
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	//create a byte array buffer to read the incoming data
	buf := make([]byte, 4096)
	// read into the buffer
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
	}

	// turn the byte array into a string
	message := string(buf)
	fmt.Println("Message received: ", message)
	// split the string into parts delimited by a " "
	parts := strings.Split(message, " ")
	// the path is the second part of the message
	path := parts[1]

	var response string
	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection:", err.Error())
			os.Exit(1)
		}
		return
	}

	if strings.HasPrefix(path, "/echo") {
		payloadParts := strings.Split(path, "/")
		// check that there are two parts to the payload
		// "echo" and  the message
		fmt.Printf("Payload parts %v length %v\n: ", payloadParts, len(payloadParts))
		if len(payloadParts) != 3 {
			response = "HTTP/1.1 400 Bad Request\r\n\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing to connection:", err.Error())
				os.Exit(1)
			}
			return
		}
		fmt.Println("Payload: ", payloadParts[2])
		response = prepareEchoResponse(payloadParts[2])
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	// write the response back to the client
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection:", err.Error())
		os.Exit(1)
	}

}
