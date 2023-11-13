package main

import (
	"log"
	"os"
)

// using os to try to begin to access ips from command line.
func main() {
	if len(os.Args) > 7 { //Shortest ip possible? (0.0.0.0)
		log.Println("Your input needs to one or multiple valid IP addresses.")
		os.Exit(400) //Http status code for wrong syntax?
	} else {
		log.Println("This is a valid IP!")
	}

	//file, err := os.Open("text.txt")
	//if err != nil {
	//log.Fatal(err)
	//}

	//data := make([]byte, 100)
	//count, err := file.Read(data)
	//if err != nil {
	//log.Fatal(err)
	//}

	//fmt.Printf("read %d bytes: %q\n", count, data[:count])
}
