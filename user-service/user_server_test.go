package main

import (
	"context"
	"testing"

	"github.com/Ali-Assar/CashWatch/authentication-service/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	// defer tdb.Exec("DROP TABLE users")

	// TODO: Create a mock userStore (you can use a real one if available)
	userStore := db.NewPostgreSQLUserStore(tdb)

	// Create a userServiceServer instance
	server := &userServiceServer{
		userStore: userStore,
	}

	// Create a sample user request
	req := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "secret",
	}

	// Call the CreateUser method
	resp, err := server.InsertUser(context.Background(), req)

	// Assertions
	assert.NoError(t, err, "Error should be nil")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.NotEmpty(t, resp.ID, "User ID should not be empty")
}
