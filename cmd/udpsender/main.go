package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	UDPSender, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalln(err)
		return
	}

	UDPCon, err := net.DialUDP("udp", nil, UDPSender)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer UDPCon.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")

		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
			return
		}
		_, err = UDPCon.Write([]byte(data))
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
