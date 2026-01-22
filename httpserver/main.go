package main

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
		body := request200()
		status := response.StatusOk
		if req.RequestLine.RequestTarget == "/urBad" {
			status = response.StatusBadRequest
			body = request500()
		} else if req.RequestLine.RequestTarget == "/urRawat" {
			status = response.StatusInternalServerError
			body = request400()
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
			target := req.RequestLine.RequestTarget
			prefix := "/httpbin"

			res, err := http.Get("https://httpbin.org" + target[len(prefix):])
			if err != nil {
				fmt.Println(err)
				status = response.StatusInternalServerError
				body = request400()
			} else {
				w.WriteStatusLine(response.StatusOk)
				h.Delete("content-length")
				h.SET("transfer-encoding", "chunked")
				h.Replace("content-type", "text/plain")
				w.WriteHeaders(*h)
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}
					//write hex
					w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					//write binary
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n\r\n"))
				return nil
			}
			defer res.Body.Close()

		}
		w.WriteStatusLine(status)
		h.Replace("content-type", "text/html")
		h.Replace("content-length", strconv.Itoa(len(body)))
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
