package request

import (
	"bytes"
	"fmt"
	"http/internal/headers"
	"io"
)

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

const (
	StateInit  = "init"
	StateDone  = "done"
	StateError = "error"
	StateParseHeader = "headers"
)

type Request struct {
	RequestLine RequestLine
	Headers *headers.Headers
	state       string
}

var ErrorMalfomredRequest = fmt.Errorf("Bad request")
var ErrorHttpVersion = fmt.Errorf("Wrong http verison")
var ErrorStateInit = fmt.Errorf("Error while state init")
var seprartor = []byte("\r\n")

func parseReqLine(s []byte) (*RequestLine, int, error) {

	i := bytes.Index(s, seprartor)
	if i == -1 {
		return nil, 0, nil
	}
	startOfLine := s[:i]
	read := i + len(seprartor)
	parts := bytes.Split(startOfLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalfomredRequest
	}

	httpParts := bytes.Split(parts[2], []byte("/"))



	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorHttpVersion
	}
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}
	return rl, read, nil
}
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		curr_data:=data[read:]
		switch r.state {
		case StateError:
			return 0, ErrorStateInit
		case StateInit:
			rl, n, err := parseReqLine(curr_data)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			r.RequestLine = *rl
			r.state = StateParseHeader

		case StateParseHeader:
			n,done,err:=r.Headers.Parse(curr_data)
			if err != nil {
				return 0, err
			}
			if n==0{
				break outer
			}
			read+=n
			if done{
				r.state=StateDone
				break outer
			}
		case StateDone:
			break outer
		default:
			panic("somehow I fuckefd up")
		}
		
	}
	return read, nil

}
func (r *Request) done() bool{
	return r.state==StateDone||r.state==StateError
}
func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
	}
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := newRequest()
	buf:=make([]byte, 1024)
	bufLn:=0
	for !req.done(){
		n,err:=reader.Read(buf[bufLn:])
		if err!=nil{
			return nil,err
		}
		bufLn+=n
		read,err:=req.parse(buf[:bufLn])
		if err!=nil{
			return nil ,err
		}
		copy(buf,buf[read:bufLn])
		bufLn-=read
	}
	return req, nil
}
