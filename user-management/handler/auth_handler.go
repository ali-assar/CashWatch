package handler

import (
	"fmt"
	"time"

	"github.com/Ali-Assar/CashWatch/types"
	"github.com/Ali-Assar/CashWatch/user-management/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: handle error
type AuthHandler struct {
	userStore db.UserStorer
}

func NewAuthHandler(userStore db.UserStorer) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("invalid credentials")
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)

}

func CreateTokenFromUser(user *types.User) string {
	expires := time.Now().Add(time.Hour * 4).Format(time.RFC3339)
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secret := "JWT_SECRET"
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign in with secret:", err)
	}
	return tokenStr
}
