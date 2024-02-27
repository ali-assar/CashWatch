package methods

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthSuccess(t *testing.T) {
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
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "secret",
	}

	_, err = server.InsertUser(context.Background(), req)
	assert.NoError(t, err, "Error should be nil")

	authParams := &pb.AuthenticateParams{
		Email:    "john@example.com",
		Password: "secret",
	}

	resp, err := server.Authenticate(context.Background(), authParams)

	assert.NoError(t, err, "Error should be nil")
	assert.NotNil(t, resp.Token)
	fmt.Println(resp.Token)
}

func TestAuthWithWrongPassword(t *testing.T) {
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
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "secret",
	}

	_, err = server.InsertUser(context.Background(), req)
	assert.NoError(t, err, "Error should be nil")

	authParams := &pb.AuthenticateParams{
		Email:    "john@example.com",
		Password: "NotSecret",
	}

	_, err = server.Authenticate(context.Background(), authParams)

	assert.Error(t, err, "Error should not be nil")

}
