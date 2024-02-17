package main

import (
	"github.com/Ali-Assar/CashWatch/types"
	"github.com/Ali-Assar/CashWatch/user-management/db"
	"github.com/gofiber/fiber"
)

type UserHandler struct {
	userStore db.UserStorer
}

func NewUserHandler(userStore db.UserStorer) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	return nil
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	return nil

}
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	return nil
}
