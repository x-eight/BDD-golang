package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	userpb "github.com/x-eight/BDD-golang/gen/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{addr}
}

func (s *Server) SendAndClient() (userpb.GreetServiceClient, error) {
	conn, err := grpc.Dial(s.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect with server %s: %w", s.Addr, err)
	}

	client := userpb.NewGreetServiceClient(conn)

	return client, nil
}

type greetServiceServer struct {
	userpb.UnimplementedGreetServiceServer
}

func doGreet(c userpb.GreetServiceClient) {
	log.Println("Unary service")
	r, err := c.Greet(context.Background(), &userpb.GreetRequest{FirstName: "saul"})

	if err != nil {
		log.Fatalf("Service failed: %v\n", err)
	}

	log.Printf("Received: %s\n", r.Result)
}

func doGreetManyTimes(c userpb.GreetServiceClient) {
	log.Println("Server streaming")

	req := &userpb.GreetRequest{FirstName: "saul"}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v\n", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}

		log.Printf("Received: %s\n", msg.Result)
	}
}

func doLongGreet(c userpb.GreetServiceClient) {
	log.Println("Client Streaming")

	reqs := []*userpb.GreetRequest{
		{FirstName: "saul"},
		{FirstName: "alonso"},
		{FirstName: "x-eight"},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling server: %v\n", err)
	}

	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server: %v\n", err)
	}

	log.Printf("Received: %s\n", res.Result)
}

func doGreetEveryone(c userpb.GreetServiceClient) {
	log.Println("Bi-direcctional service")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
	}

	requests := []*userpb.GreetRequest{
		{FirstName: "saul"},
		{FirstName: "alonso"},
		{FirstName: "x-eight"},
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			log.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error while receiving: %v\n", err)
				break
			}
			log.Printf("Received: %v\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc
}

func doGreetWithDeadline(c userpb.GreetServiceClient, timeout time.Duration) {
	log.Println("Service With Deadline")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &userpb.GreetRequest{FirstName: "saul"}
	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Code() == codes.DeadlineExceeded {
				log.Println("Deadline exceeded!")
				return
			}

			log.Fatalf("Unexpected gRPC error: %v\n", e)
		} else {
			log.Fatalf("A non gRPC error: %v\n", err)
		}
	}

	log.Printf("Received: %s\n", res.Result)
}
