package response

import (
	"fmt"
	"http/internal/headers"

	"io"
	"strconv"
)

type StatusCode int
type Writer struct{
	Writer io.Writer
}
const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func NewWriter(writer io.Writer) *Writer{
	return &Writer{Writer: writer}
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	
	h.SET("content-length", strconv.Itoa(contentLen))
	h.SET("connection", "close")
	h.SET("content-type","text-plain")
	return h
}

func (w *Writer) WriteStatusLine(status StatusCode) error{
	statusLine := []byte{}
	switch status {
	case StatusOk:
		statusLine = ([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusInternalServerError:
		statusLine = ([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	case StatusNotFound:
		statusLine = ([]byte("HTTP/1.1 404 NotFound\r\n"))
	case StatusBadRequest:
		statusLine = ([]byte("HTTP/1.1 400 BadRequest\r\n"))
	default:
		return fmt.Errorf("Unrecoginzed Error Code\r\n")
	}
	_, err := w.Writer.Write(statusLine)
	return err
}
func (w *Writer) WriteHeaders(h headers.Headers) error{
	b:=[]byte{}
	h.ForEach(func(n, v string) {
		b=fmt.Appendf(b,"%s:%s \r\n",n,v)
	})
	b=fmt.Appendf(b,"\r\n")
	_,err:=w.Writer.Write(b)
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error){
	write,err:=w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	return write,nil
}