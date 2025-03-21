package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	rAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sending message: %q", strings.TrimRight(line, "\n"))

		n, err := connection.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sent %d bytes", n)
	}
}
