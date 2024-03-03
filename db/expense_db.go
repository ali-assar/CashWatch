package db

import (
	"context"
	"database/sql"

	pb "github.com/Ali-Assar/CashWatch/types"
)

type ExpenseStorer interface {
	InsertCategory(context.Context, *pb.Category) (*pb.Category, error)
	GetCategoryById(context.Context, *pb.CategoryRequest) (*pb.Category, error)
	UpdateCategoryById(context.Context, *pb.Category) error
	DeleteCategoryById(context.Context, *pb.CategoryRequest) error

	InsertExpense(context.Context, *pb.Expense) (*pb.Expense, error)
	GetExpenseById(context.Context, *pb.ExpenseRequest) (*pb.Expense, error)
	UpdateExpenseById(context.Context, *pb.Expense) error
	DeleteExpenseById(context.Context, *pb.ExpenseRequest) error
}

type PostgreSQLExpenseStore struct {
	db *sql.DB
}

func NewPostgreSQLExpenseStore(db *sql.DB) *PostgreSQLExpenseStore {
	return &PostgreSQLExpenseStore{
		db: db,
	}
}
