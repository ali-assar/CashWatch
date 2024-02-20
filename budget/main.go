package main

import (
	"log"
	"net"

	"github.com/Ali-Assar/CashWatch/types"
	"google.golang.org/grpc"
	// Adjust import path
	// Import database or storage dependencies if needed
)

const (
	port = ":50051" // Adjust port number if needed
)

type budgetServer struct {
	types.UnimplementedBudgetServiceServer
	// Add fields or methods as needed for your actual budget service logic
	// Consider data access (database, in-memory storage), validation, authorization
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	types.RegisterBudgetServiceServer(s, &budgetServer{})

	log.Println("Budget service listening on", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
