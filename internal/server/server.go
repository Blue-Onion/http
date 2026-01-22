package server

import (

	"fmt"
	"http/internal/request"
	"http/internal/response"
	"io"
	"net"

)

type Server struct {
	Closed  bool
	Handler Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w *response.Writer, req *request.Request) *HandlerError

func createServer(handler Handler) *Server {
	return &Server{Closed: false, Handler: handler}
}
func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}

	s.Handler(responseWriter, r)

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
func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	newServer := createServer(handler)
	go runSever(newServer, listener)

	return newServer, nil
}
func (s *Server) Close() error {
	s.Closed = true
	return nil
}
