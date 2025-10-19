package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MoXcz/httpfromtcp/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer tcpListener.Close()

	fmt.Println("Listening on port 42069")
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatalln(err)
			return
		}
		fmt.Println("Connection accepted!")
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalln(err)
			return
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
	}
}
