package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"

	userclient "github.com/Ali-Assar/CashWatch/api-gateway/user-client"
	pb "github.com/Ali-Assar/CashWatch/types"
)

func main() {
	app := fiber.New()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var (
		client = pb.NewUserServiceClient(conn)
		apiV1  = app.Group("/api/v1")
	)

	apiV1.Use("/user/:email", userclient.Auth(client))
	apiV1.Use("/user/:email/*", userclient.Auth(client))

	apiV1.Post("/user", userclient.PostUser(client))
	apiV1.Get("/user/:email", userclient.GetUser(client))
	apiV1.Delete("/user/:email", userclient.DeleteUser(client))
	apiV1.Put("/user/:email", userclient.PutUser(client))

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
