package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kijudev/http-go/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v\n", req)
	}

}
