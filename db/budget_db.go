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
	query := "INSERT INTO budgets (title, amount, expireAT, setAT, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, budget.Title, budget.Amount, budget.ExpireAt, budget.SetAt, budget.UserId).Scan(&budget.ID); err != nil {
		return nil, err
	}
	return budget, nil
}

func (store *PostgreSQLBudgetStore) GetBudgetByID(ctx context.Context, id string) (*pb.Budget, error) {
	query := "SELECT title, amount, expireAT, setAT, user_id FROM budgets WHERE id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var budget pb.Budget

	if err := row.Scan(&budget.ID, &budget.Title, &budget.Amount, &budget.ExpireAt, &budget.SetAt, &budget.UserId); err != nil {
		return nil, err
	}
	return &budget, nil
}

func (store *PostgreSQLBudgetStore) UpdateBudgetByID(ctx context.Context, id string, budget *pb.Budget) error {
	query := "UPDATE budgets SET title = $1, amount = $2, expireAT = $3, setAT = $4, WHERE id = $5"
	_, err := store.db.ExecContext(ctx, query, budget.Title, budget.Amount, budget.ExpireAt, budget.SetAt, id)
	return err
}

func (store *PostgreSQLBudgetStore) DeleteBudgetByID(ctx context.Context, id string) error {
	query := "DELETE FROM budgets WHERE id = $1"
	_, err := store.db.ExecContext(ctx, query, id)
	return err
}

func (store *PostgreSQLBudgetStore) InsertIncome(ctx context.Context, income *pb.Income) (*pb.Income, error) {
	query := "INSERT INTO incomes (title, amount, receivedAt, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, income.Title, income.Amount, income.ReceivedAt, income.UserId).Scan(&income.ID); err != nil {
		return nil, err
	}
	return income, nil
}

func (store *PostgreSQLBudgetStore) GetIncomeByID(ctx context.Context, id string) (*pb.Income, error) {
	query := "SELECT title, amount, receivedAt, user_id FROM incomes WHERE id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var income pb.Income

	if err := row.Scan(&income.ID, &income.Title, &income.Amount, &income.ReceivedAt, &income.UserId); err != nil {
		return nil, err
	}
	return &income, nil
}

func (store *PostgreSQLBudgetStore) UpdateIncomeByID(ctx context.Context, id string, income *pb.Income) error {
	query := "UPDATE incomes SET title = $1, amount = $2, receivedAt = $3, WHERE id = $4"
	_, err := store.db.ExecContext(ctx, query, income.Title, income.Amount, income.ReceivedAt, id)
	return err
}

func (store *PostgreSQLBudgetStore) DeleteIncomeByID(ctx context.Context, id string) error {
	query := "DELETE FROM incomes WHERE id = $1"
	_, err := store.db.ExecContext(ctx, query, id)
	return err
}
