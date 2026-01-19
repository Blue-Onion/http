package headers

import (
	"bytes"
	"fmt"
	"strings"


)

type Headers struct {
	headers map[string]string
}

var sep = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}
func isToken(str []byte) bool{
	for _,ch:=range str{
		found:=false
		if ch>='A'&&ch<='Z'||ch>='a'&&ch<='z'||ch>='0'&&ch<='9'{
			found=true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found{
			return false
		}
	}
	return true
}
func parse(data []byte) (string, string, error) {

	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed header")
	}
	name := parts[0]
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed field Line")
	}
	value := bytes.TrimSpace(parts[1])

	return string(name), string(value), nil
}
func (h *Headers) GET(name string) string {
	return h.headers[strings.ToLower(name)]
}
func (h *Headers) SET(name, value string) {
	h.headers[strings.ToLower(name)] = value
}
func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		i := bytes.Index(data[read:], sep)
		if i == -1 {
			break
		}
		slice := data[read : read+i]

		if len(slice) == 0 {
			done = true
			read+=len(sep)
			break
		}
		name, value, err := parse(slice)
		if err != nil {
			return 0, false, err
		}
		if !isToken([]byte(name)){
			return 0,false,fmt.Errorf("Malformed Header")
		}
		h.SET(name,value)
		read += i+len(sep)
	}
	return read, done, nil
}
