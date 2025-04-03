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
		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", request.RequestLine.Method)
		fmt.Println("- Target:", request.RequestLine.RequestTarget)
		fmt.Println("- Version:", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("A connection has been closed")
	}
}
