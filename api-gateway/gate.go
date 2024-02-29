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

	app.Get("/user/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")

		resp, err := userClient.GetUserByEmail(context.Background(), &pb.UserRequest{Email: email})
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching user")
		}
		return c.JSON(resp)
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		var params pb.User
		if err := c.BodyParser(&params); err != nil {
			return c.Status(http.StatusBadRequest).SendString("failed to parse the body")
		}
		resp, err := userClient.InsertUser(context.Background(), &params)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error posting user")
		}
		return c.JSON(resp)

	})

	app.Delete("/user/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")

		_, err := userClient.DeleteUserByEmail(context.Background(), &pb.UserRequest{Email: email})
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error deleting user")
		}
		return nil
	})

	app.Put("/user/:email", func(c *fiber.Ctx) error {
		var params pb.UpdateUserRequest
		if err := c.BodyParser(&params); err != nil {
			return c.Status(http.StatusBadRequest).SendString("failed to parse the body")
		}
		params.Email = c.Params("email")

		_, err := userClient.UpdateUserByEmail(context.Background(), &params)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error updating user")
		}
		return nil
	})

	err = app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
