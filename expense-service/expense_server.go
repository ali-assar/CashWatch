package main

import (
	"log"
	"net"

	"github.com/Ali-Assar/CashWatch/db"
	expensehandler "github.com/Ali-Assar/CashWatch/expense-service/expense-handler"
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
	expenseStore := db.NewPostgreSQLExpenseStore(database)
	pb.RegisterExpenseServiceServer(server, &expensehandler.ExpenseServiceServer{
		ExpenseStore: expenseStore,
	})

	// Listen on a port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}
	log.Println("gRPC server listening on :50051")

	// Start serving requests
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
