package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// https://pkg.go.dev/os#Args
	if len(os.Args) < 2 {
		fmt.Println("Not provided an args file parameter")
		os.Exit(1)
	}
	fn := os.Args[1]

	// https://pkg.go.dev/os#Open
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("File args can't be opened")
		os.Exit(1)
	}

	// https://pkg.go.dev/io#Copy
	// https://pkg.go.dev/os@go1.18.4#pkg-variables
	// https://pkg.go.dev/os#File.Read
	io.Copy(os.Stdout, f)
}
