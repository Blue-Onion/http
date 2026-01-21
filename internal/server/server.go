package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	Closed bool
}

func createServer() *Server {
	return &Server{Closed: false}
}
func runConnection(s *Server, conn io.ReadWriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()

}
func runSever(s *Server, listener net.Listener) error {

	go func() {
		for {
			if s.Closed {
				return
			}
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go runConnection(s, conn)
		}

	}()
	return nil
}
func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	newServer := createServer()
	go runSever(newServer, listener)

	return newServer, nil
}
func (s *Server) Close() error {
	s.Closed = true
	return nil
}
