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
	}

	//else {
	//log.Println("This is a valid IP!")
	//}

	ipport := os.Args[1]
	ipportnext := os.Args[2]
	log.Println(ipport)
	log.Println(ipportnext)

	//give each node a corresponding id.
	//that id needs to be something the other nodes will be able to become aware of.
	//needs to be a critical section somewhere that each node needs to access at some point.
	//Use a token ring to pass information in a circle, such that every node gets a turn.
	//Use unlock and lock, to keep the information to only one node at a time.

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
