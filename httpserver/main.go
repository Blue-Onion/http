package main

import (
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"io"

	"log"
	"os"
	"os/signal"
	"syscall"


)


const port = 42069



func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError{
		if req.RequestLine.RequestTarget=="/myBad"{
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:"vfrf",
			}
		}else{
			w.Write([]byte("All good madafaka \n"))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}