package userclient

import (
	"net/http"
	"time"

	pb "github.com/Ali-Assar/CashWatch/types"
	userhandler "github.com/Ali-Assar/CashWatch/user-service/user-handler"
	"github.com/gofiber/fiber/v2"
)

func Auth(userClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokens, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok || len(tokens) == 0 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		token := tokens[0]
		claims, err := userhandler.ValidateToken(token)
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

		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func GetUser(userClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		email := c.Params("email")
		if user.Email != email {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		resp, err := userClient.GetUserByEmail(c.Context(), &pb.UserRequest{Email: email})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user"})
		}
		return c.JSON(resp)
	}
}

func PostUser(userClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params pb.User
		if err := c.BodyParser(&params); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse the body"})
		}
		resp, err := userClient.InsertUser(c.Context(), &params)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error posting user"})
		}
		return c.JSON(resp)
	}
}

func DeleteUser(userClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		email := c.Params("email")
		if user.Email != email {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		_, err := userClient.DeleteUserByEmail(c.Context(), &pb.UserRequest{Email: email})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting user"})
		}
		return c.SendStatus(http.StatusOK)
	}
}

func PutUser(userClient pb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params pb.UpdateUserRequest
		if err := c.BodyParser(&params); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse the body"})
		}

		user, ok := c.Context().UserValue("user").(*pb.User)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		params.Email = c.Params("email")
		if user.Email != params.Email {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		_, err := userClient.UpdateUserByEmail(c.Context(), &params)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating user"})
		}
		return c.SendStatus(http.StatusOK)
	}
}
