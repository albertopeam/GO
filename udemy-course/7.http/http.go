package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWritter struct{}

func (logWritter) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func main() {
	fmt.Println("***************************")
	fmt.Println("********* MANUAL **********")
	fmt.Println("***************************")
	manual()
	fmt.Println("")
	fmt.Println("***************************")
	fmt.Println("********* AUTO **********")
	fmt.Println("***************************")
	auto()
	fmt.Println("")
	fmt.Println("***************************")
	fmt.Println("********* WRITER **********")
	fmt.Println("***************************")
	writer()
}

// manual read body
func manual() {
	// https://pkg.go.dev/net/http@go1.18.4#pkg-index
	resp, err := http.Get("https://google.com")
	if err != nil {
		fmt.Println("Error fetching: ", err)
		os.Exit(1)
	}
	// https://pkg.go.dev/net/http@go1.18.4#Response
	// https://pkg.go.dev/io#ReadCloser - implements two interfaces, Reader and Closer
	bs := make([]byte, 99999)
	size, err := resp.Body.Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading: ", err)
		os.Exit(1)
	}
	fmt.Println("Bytes read: ", size)
	fmt.Println(string(bs))
}

// automatic read body
func auto() {
	// https://pkg.go.dev/net/http@go1.18.4#pkg-index
	resp, err := http.Get("https://google.com")
	if err != nil {
		fmt.Println("Error fetching: ", err)
		os.Exit(1)
	}
	// https://pkg.go.dev/io#Writer
	// https://pkg.go.dev/io#Copy
	io.Copy(os.Stdout, resp.Body)
}

// custom read body with a Writer
func writer() {
	// https: //pkg.go.dev/io@go1.18.4#Writer
	// https://pkg.go.dev/net/http@go1.18.4#pkg-index
	resp, err := http.Get("https://google.com")
	if err != nil {
		fmt.Println("Error fetching: ", err)
		os.Exit(1)
	}
	// https://pkg.go.dev/io#Writer
	// https://pkg.go.dev/io#Copy

	lw := logWritter{}
	io.Copy(lw, resp.Body)
}
