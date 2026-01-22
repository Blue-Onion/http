package main

import (
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const port = 42069

func request400() []byte {

	return []byte(`<html>
	<head>
	  <title>400 Bad Request</title>
	</head>
	<body>
	  <h1>Bad Request</h1>
	  <p>Your request honestly kinda sucked.</p>
	</body>
	</html>`)
}
func request500() []byte {

	return []byte(`<html>
	<head>
	  <title>500 Internal Server Error</title>
	</head>
	<body>
	  <h1>Internal Server Error</h1>
	  <p>Okay, you know what? This one is on me.</p>
	</body>
  </html>`)
}
func request200() []byte {

	return []byte(`<html>
	<head>
	  <title>200 OK</title>
	</head>
	<body>
	  <h1>Success!</h1>
	  <p>Your request was an absolute banger.</p>
	</body>
  </html>`)
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) *server.HandlerError {
		h := response.GetDefaultHeaders(0)
		body:=request200()
		status:=response.StatusOk
		switch req.RequestLine.RequestTarget{
		case "/urBad":
			status=response.StatusBadRequest
			body=request500()
		case  "/urRawat":
			status=response.StatusBadRequest
			body=request500()
		}
	
		w.WriteStatusLine(status)
		h.Replace("content-length",strconv.Itoa(len(body)))
		h.Replace("content-type","text/html")
		w.WriteHeaders(*h)
		w.WriteBody(body)
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
