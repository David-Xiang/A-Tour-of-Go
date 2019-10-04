package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (reader rot13Reader) Read(b []byte) (int, error) {
	n, err := reader.r.Read(b)
	if n > 0 {
		for i:=0; i < n; i++ {
			if (b[i] <= 'M' && b[i] >= 'A') || 
				(b[i] <= 'm' && b[i] >= 'a') {
				b[i] += 13
			} else if (b[i] <= 'Z' && b[i] >= 'N') || 
				(b[i] <= 'z' && b[i] >= 'n'){
				b[i] -= 13
			}
		}
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}