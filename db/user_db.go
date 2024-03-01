package db

import (
	"context"
	"database/sql"

	pb "github.com/Ali-Assar/CashWatch/types"
)

type UserStorer interface {
	GetUserByEmail(context.Context, string) (*pb.User, error)
	GetUserByID(context.Context, any) (*pb.User, error)
	GetUsers(context.Context) ([]*pb.User, error)
	InsertUser(context.Context, *pb.User) (*pb.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *pb.UpdateUserParams) error
	UpdateUserByEmail(context.Context, string, *pb.UpdateUserRequest) error
}

type PostgreSQLUserStore struct {
	db *sql.DB
}

func NewPostgreSQLUserStore(db *sql.DB) *PostgreSQLUserStore {
	return &PostgreSQLUserStore{
		db: db,
	}
}

func (store *PostgreSQLUserStore) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	query := "SELECT id,firstName, lastName, email, encryptedPassword FROM users where email = $1"
	row := store.db.QueryRowContext(ctx, query, email)
	var user pb.User

	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *PostgreSQLUserStore) GetUserByID(ctx context.Context, id any) (*pb.User, error) {
	query := "SELECT id, firstName, lastName, email FROM users where id = $1"
	row := store.db.QueryRowContext(ctx, query, id)
	var user pb.User

	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *PostgreSQLUserStore) GetUsers(ctx context.Context) ([]*pb.User, error) {
	query := "SELECT id, firstName, lastName, email FROM users"
	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (store *PostgreSQLUserStore) InsertUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	query := "INSERT INTO users(firstName, lastName, email, encryptedPassword) VALUES($1, $2, $3, $4) RETURNING id"
	if err := store.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

func (store *PostgreSQLUserStore) DeleteUser(ctx context.Context, email string) error {
	query := "DELETE FROM users WHERE email = $1"
	_, err := store.db.ExecContext(ctx, query, email)
	return err
}

func (store *PostgreSQLUserStore) UpdateUser(ctx context.Context, email string, user *pb.UpdateUserParams) error {
	query := "UPDATE users SET firstName = $1, lastName = $2 WHERE email = $3"
	_, err := store.db.ExecContext(ctx, query, user.FirstName, user.LastName, email)
	return err
}

func (store *PostgreSQLUserStore) UpdateUserByEmail(ctx context.Context, email string, user *pb.UpdateUserRequest) error {
	query := "UPDATE users SET firstName = $1, lastName = $2 WHERE email = $3"
	_, err := store.db.ExecContext(ctx, query, user.FirstName, user.LastName, email)
	return err
}
