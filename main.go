package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	ping "github.com/eepyITU/Distributed-Mutual-Exclusion/grpc"
	"google.golang.org/grpc"
)

func main() {
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

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
	}()

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

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

	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//p.sendPingToNeighbor()
	//}
	for {
		rand := rand.Intn(100)
		if rand > 90 {
			p.sendPingToNeighbor()
		} else {
			p.ImportantWork()
		}
	}
	//istedet for at bede om request og sender dét rundt
	//istedet tjek om p har token, og om den skal sendes videre

}

type Peer struct {
	ping.UnimplementedPingServer
	id            int32
	amountOfPings map[int32]int32
	clients       map[int32]ping.PingClient
	ctx           context.Context
	state         bool
}

func (p *Peer) Ping(ctx context.Context, req *ping.Request) (*ping.Reply, error) {
	id := req.Id
	p.amountOfPings[id] += 1

	//er den selv i critical section
	requestToken := req.RequestToken
	if req.RequestId == p.id {
		//hvem har sendt requesten originalt (propagerer id'et videre)
		if requestToken == 1 {
			p.enterCriticalSection()
			return &ping.Reply{
				Id:    p.id,
				Reply: 1,
			}, nil
		} else {
			log.Printf("Peer %d was denied access to the critical section.", p.id)
			return &ping.Reply{
				Id:    p.id,
				Reply: 1,
			}, nil
		}
	} else {
		if p.state {
			requestToken = 0
		}

		//sende videre hvis det ikke er den selv
		p.propagatePingToNeighbor(req.RequestId, requestToken)
		return &ping.Reply{
			Id:    p.id,
			Reply: 1,
		}, nil

	}

}

//This func goes from 'SendPingToAll' to 'SendPingToNeighbor' instead.
//Since it needs only ping to the next port in the sequence.
//sendTokenToNeighbor (--p.id)
//hvis nuværende peer har token, state true
//ellers false (tjekker om den har token når den enter critical section.)
//efter critical section, skal den sendes videre
//at enter critical section skal du have token,
//kan ikke enter uden token

func (p *Peer) sendPingToNeighbor() {
	request := &ping.Request{Id: p.id, RequestToken: 1, RequestId: p.id}

	nextPort := (p.id + 1) % int32(len(p.clients))

	client := p.clients[nextPort]

	reply, err := client.Ping(p.ctx, request)
	if err != nil {
		fmt.Println("something went wrong")
	}
	fmt.Printf("Got reply from id %v: %v\n", reply.Id, reply.Reply)

	//p.enterCriticalSection()

}

func (p *Peer) propagatePingToNeighbor(requestId int32, requestToken int32) {
	request := &ping.Request{Id: p.id, RequestToken: requestToken, RequestId: int32(requestId)}

	nextPort := (p.id + 1) % int32(len(p.clients))

	client := p.clients[nextPort]

	reply, err := client.Ping(p.ctx, request)
	if err != nil {
		fmt.Println("something went wrong")
	}
	fmt.Printf("Got reply from id %v: %v\n", nextPort, reply.Reply)

	//p.enterCriticalSection()

}

// trying to make a function that makes the peer access the critical section.
func (p *Peer) enterCriticalSection() {
	log.Printf("Peer %d is now in the critical section.\n", p.id)

	time.Sleep(3)

	log.Printf("Peer %d is now done with the critical section.\n", p.id)
}

func (p *Peer) ImportantWork() {
	log.Printf("Peer %d is starting some very important work.\n", p.id)

	time.Sleep(3)

	log.Printf("Peer %d is finishing some very important work.\n", p.id)

}
