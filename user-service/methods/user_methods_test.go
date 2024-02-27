package methods

import (
	"context"
	"testing"

	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

	// TODO: Create a mock userStore (you can use a real one if available)
	userStore := db.NewPostgreSQLUserStore(tdb)

	server := &UserServiceServer{
		UserStore: userStore,
	}
	req := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "secret",
	}

	resp, err := server.InsertUser(context.Background(), req)

	assert.NoError(t, err, "Error should be nil")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NotEmpty(t, resp.ID, "User ID should not be empty")
}

func TestDeleteUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

	userStore := db.NewPostgreSQLUserStore(tdb)

	server := &UserServiceServer{
		UserStore: userStore,
	}
	req := &pb.User{
		FirstName: "james",
		LastName:  "foo",
		Email:     "james@foo.com",
		Password:  "secret",
	}

	_, err = server.InsertUser(context.Background(), req)
	assert.NoError(t, err, "Error should be nil")

	email := &pb.UserRequest{Email: "james@foo.com"}
	_, err = server.DeleteUserByEmail(context.Background(), email)

	assert.NoError(t, err, "Error should be nil")
}

func TestGetUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}

	defer tdb.Exec("DROP TABLE users")

	userStore := db.NewPostgreSQLUserStore(tdb)

	server := &UserServiceServer{
		UserStore: userStore,
	}
	req := &pb.User{
		FirstName: "james",
		LastName:  "foo",
		Email:     "james@foo.com",
		Password:  "secret",
	}

	_, err = server.InsertUser(context.Background(), req)
	assert.NoError(t, err, "Error should be nil")

	email := &pb.UserRequest{Email: "james@foo.com"}
	user, err := server.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, req.FirstName, user.FirstName)
	assert.Equal(t, req.LastName, user.LastName)
	assert.Equal(t, req.Email, user.Email)
	assert.Empty(t, user.Password)
}

func TestUpdateUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}

	defer tdb.Exec("DROP TABLE users")

	userStore := db.NewPostgreSQLUserStore(tdb)

	server := &UserServiceServer{
		UserStore: userStore,
	}
	req := &pb.User{
		FirstName: "james",
		LastName:  "foo",
		Email:     "james@foo.com",
		Password:  "secret",
	}

	_, err = server.InsertUser(context.Background(), req)
	assert.NoError(t, err, "Error should be nil")

	updateParam := &pb.UpdateUserRequest{
		FirstName: "bar",
		LastName:  "baz",
		Email:     "james@foo.com",
	}
	_, err = server.UpdateUserByEmail(context.Background(), updateParam)
	assert.NoError(t, err, "Error should be nil")

	email := &pb.UserRequest{Email: "james@foo.com"}
	user, err := server.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, updateParam.FirstName, user.FirstName)
	assert.Equal(t, updateParam.LastName, user.LastName)

	// Test with an invalid email (non-existent user)
	updateInvalidParam := &pb.UpdateUserRequest{
		FirstName: "james",
		LastName:  "foo",
		Email:     "nonexistent@example.com", // Invalid email
	}
	server.UpdateUserByEmail(context.Background(), updateInvalidParam)

	user, err = server.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err, "Error should be nil")
	assert.NotEqual(t, updateInvalidParam.FirstName, user.FirstName)
	assert.NotEqual(t, updateInvalidParam.LastName, user.LastName)

}
