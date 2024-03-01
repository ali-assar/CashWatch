package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"

	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/Ali-Assar/CashWatch/user-service/methods"
)

func main() {
	app := fiber.New()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	// Middleware for JWT authentication
	app.Use(func(c *fiber.Ctx) error {
		tokens, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok || len(tokens) == 0 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		token := tokens[0]
		claims, err := methods.ValidateToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		expiresStr, ok := claims["expires"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		expires, err := time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials", "details": err.Error()})
		}

		if time.Now().After(expires) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		email := claims["email"]

		strEmail, _ := email.(string)
		user, err := userClient.GetUserByEmail(c.Context(), &pb.UserRequest{Email: strEmail})
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		//set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	})

	app.Get("/user/:email", func(c *fiber.Ctx) error {
		// Extract user information from context (already set by middleware)
		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		email := c.Params("email")

		// Check if the authenticated user has permission to access this email
		if user.Email != email {
			return fmt.Errorf("unauthorized")
		}

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
		// Extract user information from context (already set by middleware)
		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		email := c.Params("email")
		if user.Email != email {
			return fmt.Errorf("unauthorized")
		}

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
		// Extract user information from context (already set by middleware)
		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return fmt.Errorf("unauthorized")
		}

		params.Email = c.Params("email")
		if user.Email != params.Email {
			return fmt.Errorf("unauthorized")
		}

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
