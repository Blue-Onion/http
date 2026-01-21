package server

import (
	"bytes"
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"io"
	"net"
	"strconv"
)

type Server struct {
	Closed  bool
	Handler Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w io.Writer, req *request.Request) *HandlerError

func createServer(handler Handler) *Server {
	return &Server{Closed: false, Handler: handler}
}
func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	header := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, header)
		
		return
	}
	writer := bytes.NewBuffer([]byte{})
	herr:=s.Handler(writer,r)
	if herr!=nil{
		response.WriteStatusLine(conn, herr.StatusCode)
		response.WriteHeaders(conn, header)
		conn.Write([]byte(herr.Message))
		return

	}
	body:=writer.Bytes()
	header.Replace("content-length",strconv.Itoa(len(body)))

	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, header)
	conn.Write(body)

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
