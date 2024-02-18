package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Ali-Assar/CashWatch/types"
	"github.com/Ali-Assar/CashWatch/user-management/db"
	"github.com/gofiber/fiber/v2"
)

func createUserAndPost(t *testing.T, app *fiber.App, userHandler *UserHandler, params types.CreateUserParams) *types.User {
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}
	return &user
}

func TestPostUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Close()
	userStore := db.NewPostgreSQLUserStore(tdb)
	app := fiber.New()

	userHandler := NewUserHandler(userStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foo@bar.com",
		Password:  "foobar123",
	}

	user := createUserAndPost(t, app, userHandler, params)

	//TODO: do not return EncryptedPassword
	// if len(user.EncryptedPassword) > 0 {
	// 	t.Errorf("EncryptedPassword should not included in json response")
	// }

	// Ensure the response status is OK (200)
	if user == nil {
		t.Fatal("user is nil")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected the first name be %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected the last name be %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected the email be %s but got %s", params.Email, user.Email)
	}
}
