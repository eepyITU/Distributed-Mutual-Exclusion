package main

import (
	"fmt"
	"log"
	"os"
)

// using os to try to begin to access ips from command line.
func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("read %d bytes: %q\n", count, data[:count])
}
