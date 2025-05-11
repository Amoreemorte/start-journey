package main

import (
	"log"
	"net"
	"net/rpc"
)

type Server struct {
	serverMessages []string
}

func (s *Server) Send(request []string, responce *string) error {
	s.serverMessages = append(s.serverMessages, request...)

	return nil
}

func (s *Server) Messages(_ struct{}, responce *[]string) error {
	*responce = s.serverMessages
	return nil
}

func main() {
	var RPCServer Server
	if err := rpc.Register(&RPCServer); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Print(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}

		go rpc.ServeConn(conn)
	}
}
