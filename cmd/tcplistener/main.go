package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("A connection has been accepted")
		requestLine, err := request.RequestFromReader(connection) 
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println("Request line:")
    fmt.Println("- Method:", requestLine.RequestLine.Method)
    fmt.Println("- Target:", requestLine.RequestLine.RequestTarget)
    fmt.Println("- Version:", requestLine.RequestLine.HttpVersion)
		fmt.Println("A connection has been closed")
	}
}
