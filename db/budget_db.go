package db

import (
	"context"
	"database/sql"

	pb "github.com/Ali-Assar/CashWatch/types"
)

type BudgetStorer interface {
	InsertBudget(context.Context, *pb.Budget) (*pb.Budget, error)
	GetBudgetByID(context.Context, string) (*pb.Budget, error)
	UpdateBudgetByID(context.Context, string, *pb.Budget) error
	DeleteBudgetByID(context.Context, string) error

	InsertIncome(context.Context, *pb.Income) (*pb.Income, error)
	GetIncomeByID(context.Context, string) (*pb.Income, error)
	UpdateIncomeByID(context.Context, string, *pb.Income) error
	DeleteIncomeByID(context.Context, string) error
}

type PostgreSQLBudgetStore struct {
	db *sql.DB
}

func NewPostgreSQLBudgetStore(db *sql.DB) *PostgreSQLBudgetStore {
	return &PostgreSQLBudgetStore{
		db: db,
	}
}

func (store *PostgreSQLBudgetStore) InsertBudget(ctx context.Context, budget *pb.Budget) (*pb.Budget, error) {
	query := "INSERT INTO budgets (title, amount, expireAT, setAT, user_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, budget.Title, Budget.Amount, Budget.ExpireAT, Budget.SetAt, budget.Userid).Scan(&Budget.ID); err != nil {
		return nil, err
	}
	return Budget, nil
}
