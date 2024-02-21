package handler

import (
	"strconv"

	"github.com/Ali-Assar/CashWatch/authentication-service/db"
	"github.com/Ali-Assar/CashWatch/types"
	"github.com/gofiber/fiber/v2"
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
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	// if errors := params.ValidateUpdate; len(errors) > 0 {
	// 	return c.JSON(errors)
	// }
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	if h.userStore.UpdateUser(c.Context(), id, &params); err != nil {
		return err
	}
	return c.JSON(map[string]int{"updated": id})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFormParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)

}
func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")
	if err := h.userStore.DeleteUser(c.Context(), email); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": email})
}
