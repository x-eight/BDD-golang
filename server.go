package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	userpb "github.com/x-eight/BDD-golang/gen/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var greetWithDeadlineTime time.Duration = 1 * time.Second

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{addr}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.Addr, err)
	}

	server := grpc.NewServer()
	userpb.RegisterGreetServiceServer(server, &greetServiceServer{})
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

type greetServiceServer struct {
	userpb.UnimplementedGreetServiceServer
}

func (s *greetServiceServer) Greet(ctx context.Context, req *userpb.GreetRequest) (*userpb.GreetResponse, error) {
	firstName := req.GetFirstName()

	log.Printf(firstName)

	return &userpb.GreetResponse{Result: fmt.Sprintf("Hello %s", firstName)}, nil
}

func (s *greetServiceServer) GreetManyTimes(req *userpb.GreetRequest, stream userpb.GreetService_GreetManyTimesServer) error {
	log.Printf("Server streaming with %v\n", req)
	firstName := req.GetFirstName()
	for i := 0; i < 10; i++ {
		res := fmt.Sprintf("Hello %s, number %d", firstName, i)

		stream.Send(&userpb.GreetResponse{
			Result: res,
		})
	}

	return nil
}

func (s *greetServiceServer) LongGreet(stream userpb.GreetService_LongGreetServer) error {
	log.Println("Client Streaming")

	res := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&userpb.GreetResponse{
				Result: res,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		log.Printf("Receiving req: %v\n", req)
		res += "Hello " + req.FirstName + "!\n"
	}

	return nil
}

func (s *greetServiceServer) GreetEveryone(stream userpb.GreetService_GreetEveryoneServer) error {
	log.Println("Bi-direcctional service")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		res := "Hello " + req.FirstName + "!"

		err = stream.Send(&userpb.GreetResponse{
			Result: res,
		})

		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
		}
	}

	return nil
}

func (s *greetServiceServer) GreetWithDeadline(ctx context.Context, req *userpb.GreetRequest) (*userpb.GreetResponse, error) {
	log.Printf("Service With Deadline with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(greetWithDeadlineTime)
	}

	firstName := req.GetFirstName()

	return &userpb.GreetResponse{Result: "Hello " + firstName}, nil

}
