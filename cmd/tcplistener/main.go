package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatalln(err)
			return
		}
		fmt.Println("Connection accepted!")
		linesCh := getLinesChannel(conn)
		for line := range linesCh {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)
	go func() {
		defer f.Close()
		defer close(linesCh)
		data := make([]byte, 8)
		bytesRead, err := f.Read(data)
		str := ""
		for err != io.EOF {
			str = str + string(data[0:bytesRead])
			bytesRead, err = f.Read(data)
		}
		listStr := strings.SplitSeq(str, "\n")
		for str := range listStr {
			if str != "" {
				linesCh <- str
			}
		}
	}()
	return linesCh
}
