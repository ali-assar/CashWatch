package db

import (
	"context"
	"database/sql"

	pb "github.com/Ali-Assar/CashWatch/types"
)

type ExpenseStorer interface {
	InsertCategory(context.Context, *pb.Category) (*pb.Category, error)
	GetCategoryById(context.Context, string) (*pb.Category, error)
	UpdateCategoryById(context.Context, string, *pb.Category) error
	DeleteCategoryById(context.Context, string) error

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

func (store *PostgreSQLExpenseStore) InsertCategory(ctx context.Context, category *pb.Category) (*pb.Category, error) {
	query := "INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, category.Name, category.UserId).Scan(&category.ID); err != nil {
		return nil, err
	}
	return category, nil
}

func (store *PostgreSQLExpenseStore) GetCategoryById(ctx context.Context, id string) (*pb.Category, error) {
	query := "SELECT id, name, user_id FROM categories WHERE id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var category pb.Category

	if err := row.Scan(&category.ID, &category.Name, &category.UserId); err != nil {
		return nil, err
	}
	return &category, nil
}

func (store *PostgreSQLExpenseStore) UpdateCategoryById(ctx context.Context, id string, category *pb.Category) error {
	query := "UPDATE categories SET name = $1 WHERE id = $2"
	_, err := store.db.ExecContext(ctx, query, category.Name, id)
	return err
}

func (store *PostgreSQLExpenseStore) DeleteCategoryById(ctx context.Context, id *pb.CategoryRequest) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := store.db.ExecContext(ctx, query, id)
	return err
}

func (store *PostgreSQLExpenseStore) InsertExpense(ctx context.Context, expense *pb.Expense) (*pb.Expense, error) {
	query := "INSERT INTO expenses (user_id, description, amount) VALUES ($1, $2, $3) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, expense.UserId, expense.Description, expense.Amount).Scan(&expense.Id); err != nil {
		return nil, err
	}
	return expense, nil
}

func (store *PostgreSQLExpenseStore) GetExpenseById(ctx context.Context, id *pb.ExpenseRequest) (*pb.Expense, error) {
	query := "SELECT id, user_id, description, amount, category_id FROM expenses WHERE id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var expense pb.Expense

	if err := row.Scan(&expense.UserId, expense.Description, expense.Amount, expense.CategoryId); err != nil {
		return nil, err
	}
	return &expense, nil
}

func (store *PostgreSQLExpenseStore) UpdateExpenseById(ctx context.Context, id string, expense *pb.Expense) error {
	query := "UPDATE expenses SET description = $1, amount = $2, category_id = $3  WHERE id = $4"
	_, err := store.db.ExecContext(ctx, query, expense.Description, expense.Amount, expense.CategoryId, id)
	return err
}

func (store *PostgreSQLExpenseStore) DeleteExpenseById(ctx context.Context, id *pb.ExpenseRequest) error {
	query := "DELETE FROM expenses WHERE id = $1"
	_, err := store.db.ExecContext(ctx, query, id)
	return err
}
