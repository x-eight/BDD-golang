package main

import (
	"context"
	"fmt"
	"net"

	userpb "github.com/x-eight/BDD-golang/gen/proto/user/v1"
	"google.golang.org/grpc"
)

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
	return &userpb.GreetResponse{Result: fmt.Sprintf("Hello %s", firstName)}, nil
}
