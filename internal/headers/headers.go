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
func (h *Headers) ForEach(cb func(n, v string)) {
	for i, v := range h.headers {
		cb(i, v)
	}
}
func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
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
func (h *Headers) GET(name string) (string, bool) {
	if h == nil || h.headers == nil {
		return "", false
	}
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}
func (h *Headers) SET(name, value string) {

	name = strings.ToLower(name)
	newValue := value
	val, ok := h.headers[name]
	if ok {
		h.headers[name] = val + ", " + value
	} else {
		h.headers[name] = newValue
	}
}
func (h *Headers) Replace(name, value string) {

	name = strings.ToLower(name)
	h.headers[name] = value
}
func (h *Headers) Delete(name string) {
	name = strings.ToLower(name)
	delete(h.headers, name)
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
			read += len(sep)
			break
		}
		name, value, err := parse(slice)
		if err != nil {
			return 0, false, err
		}
		if !isToken([]byte(name)) {

			return 0, false, fmt.Errorf("Malformed Header")
		}
		h.SET(name, value)
		read += i + len(sep)
	}
	return read, done, nil
}
