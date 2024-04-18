package main

import (
	"log"
	"net"

	budgethandler "github.com/Ali-Assar/CashWatch/budget-service/budget-handlers"
	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"google.golang.org/grpc"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	server := grpc.NewServer()
	budgetStore := db.NewPostgreSQLBudgetStore(database)
	pb.RegisterBudgetServiceServer(server, &budgethandler.BudgetServiceServer{
		BudgetStore: budgetStore,
	})

	// Listen on a port
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}
	log.Println("gRPC server listening on :50051")

	// Start serving requests
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
