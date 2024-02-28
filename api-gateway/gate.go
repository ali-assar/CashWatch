package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"

	pb "github.com/Ali-Assar/CashWatch/types"
)

func main() {
	app := fiber.New()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		resp, err := userClient.GetUserByEmail(context.Background(), &pb.UserRequest{Email: id})
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching user")
		}
		return c.JSON(resp)
	})

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
