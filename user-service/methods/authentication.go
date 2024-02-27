package methods

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/golang-jwt/jwt/v5"
)

func (h *UserServiceServer) Authenticate(ctx context.Context, req *pb.AuthenticateParams) (*pb.AuthenticateResp, error) {
	user, err := h.UserStore.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	if !pb.IsPasswordValid(user.Password, req.GetPassword()) {
		return nil, fmt.Errorf("invalid credentials")
	}

	resp := pb.AuthenticateResp{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return &resp, nil
}

func CreateTokenFromUser(user *pb.User) string {
	expires := time.Now().Add(time.Hour * 4).Format(time.RFC3339)
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "JWTSECRET"
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign in with secret:", err)
	}
	return tokenStr
}
