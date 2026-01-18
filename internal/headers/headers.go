package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var sep = []byte("\r\n")

func NewHeaders() Headers {

	return map[string]string{}
}
func parseHeader(fldLine []byte) (string, string, error) {
	fmt.Println("Gone here")
	parts := bytes.SplitN(fldLine, []byte(":"), 2)
	if len(parts) != 2 {

		return "", "", fmt.Errorf("Bitch wrong header")
	}
	fmt.Println(string(fldLine))
	key := parts[0]
	fmt.Println(string(key))
	if bytes.HasSuffix(key, []byte(" ")) {

		return "", "", fmt.Errorf("Naughty boi wrong field Line")
	}
	value := bytes.TrimSpace(parts[1])

	return string(key), string(value), nil

}
func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		// Host: localhost:42067\r\n\r\n
		// 0
		i := bytes.Index(data[read:], sep)
		if i == -1 {
			break
		}
		//Header Empty
		if len(data[read : i+read])==0{
			done = true
			break
		}
		//Read
		
		key, value, err := parseHeader(data[read : i+read])
		if err != nil {
			return 0, false, err
		}
		read += i + len(sep)
		h[key] = value

	}
	return read, done, nil
}
