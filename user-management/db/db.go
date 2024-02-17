package db

import (
	"context"
	"database/sql"

	"github.com/Ali-Assar/CashWatch/types"
)

var DBNAME = "user-management"

type UserStorer interface {
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUserByID(context.Context, int) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(*types.User) error
}

type PostgreSQLUserStore struct {
	db *sql.DB
}

func NewPostgreSQLUserStore(db *sql.DB) *PostgreSQLUserStore {
	return &PostgreSQLUserStore{
		db: db,
	}
}

func (store *PostgreSQLUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	query := "SELECT id, firstName, lastName , email FROM users where email = $1"
	row := store.db.QueryRowContext(ctx, query, email)
	var user types.User

	if err := row.Scan(&user.ID, &user.Email, &user.EncryptedPassword); err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *PostgreSQLUserStore) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	query := "SELECT id, firstName, lastName, email FROM users where id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var user types.User

	if err := row.Scan(&user.ID, &user.Email, &user.EncryptedPassword); err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *PostgreSQLUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	query := "SELECT id, firstName, lastName, email FROM users"
	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Email, &user.EncryptedPassword); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (store *PostgreSQLUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	query := "INSERT INTO users(firstName, lastName, email, encryptedPassword) VALUES($1,$2) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.EncryptedPassword).Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (store *PostgreSQLUserStore) DeleteUser(ctx context.Context, email string) error {
	query := "DELETE FROM users WHERE email = $1"
	_, err := store.db.ExecContext(ctx, query, email)
	return err
}

func (store *PostgreSQLUserStore) UpdateUser(ctx context.Context, user *types.User) error {
	query := "UPDATE users SET firstName = $1, lastName = $2, email = $3, password = $4 WHERE id = $5"
	_, err := store.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.EncryptedPassword, user.ID)
	return err
}
