package db

import (
	"context"

	"github.com/Ali-Assar/CashWatch/types"
)

var DBNAME = "user-management"

type UserStore interface {
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(*types.User) error
}
