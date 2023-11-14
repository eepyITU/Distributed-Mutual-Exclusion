package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	ping "github.com/eepyITU/Distributed-Mutual-Exclusion/grpc"
	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &Peer{
		id:            ownPort,
		amountOfPings: make(map[int32]int32),
		clients:       make(map[int32]ping.PingClient),
		ctx:           ctx,
	}

	//NOW: A token ring
	//1: Create a critical section that nodes want to access
	//2: The token ring in question is already here? (With dialing), specifically for 3 nodes (HARDCODED)
	//3: Make sure that the first node (1) starts the ring
	//4: 0 sends data to 1
	//5: 1 sends data to 2
	//6: 2 sends data to 0

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	ping.RegisterPingServer(grpcServer, p)

	go func() {
		//In here is where we need to have the critical section token.
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
		log.Println("This is a critical section token.")
	}()

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

		p.enterCriticalSection()

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := ping.NewPingClient(conn)
		p.clients[port] = c

	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.sendPingToNeighbor()
	}
}

type Peer struct {
	ping.UnimplementedPingServer
	id            int32
	amountOfPings map[int32]int32
	clients       map[int32]ping.PingClient
	ctx           context.Context
	lock          sync.Mutex
}

func (p *Peer) Ping(ctx context.Context, req *ping.Request) (*ping.Reply, error) {
	id := req.Id
	p.amountOfPings[id] += 1

	rep := &ping.Reply{Amount: p.amountOfPings[id]}
	return rep, nil
}

//This func goes from 'SendPingToAll' to 'SendPingToNeighbor' instead.
//Since it needs only ping to the next port in the sequence.

func (p *Peer) sendPingToNeighbor() {
	fmt.Printf("sendPingToNeighbor() is starting.")
	request := &ping.Request{Id: p.id}

	nextPort := (p.id + 1) % int32(len(p.clients))

	client := p.clients[nextPort]

	reply, err := client.Ping(p.ctx, request)
	if err != nil {
		fmt.Println("something went wrong")
	}
	fmt.Printf("Got reply from id %v: %v\n", nextPort, reply.Amount)

	//p.enterCriticalSection()

}

// trying to make a function that makes the peer access the critical section.
func (p *Peer) enterCriticalSection() {
	p.lock.Lock()
	defer p.lock.Unlock()
	//defer func
	log.Printf("Peer %d is now in the critical section.\n", p.id)

	log.Printf("Peer %d is now done with the critical section.\n", p.id)
}
