package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(linesChan)
		line := ""
		for {
			data := make([]byte, 8)
			count, err := f.Read(data)
			if err == io.EOF {
				break
			} else {
				str := string(data[:count])
				parts := strings.Split(str, "\n")
				for i := range len(parts) - 1 {
					linesChan <- line + parts[i]
					line = ""
				}

				line += parts[len(parts)-1]
			}
		}
		if line != "" {
			linesChan <- line
		}
	}()
	return linesChan
}

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
		linesChan := getLinesChannel(connection)
		for line := range linesChan {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("A connection has been closed")
	}
}
