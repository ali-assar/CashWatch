package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Ali-Assar/CashWatch/authentication-service/db"
	"github.com/Ali-Assar/CashWatch/types"
	"github.com/gofiber/fiber/v2"
)

func TestAuthSuccess(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

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
	_ = createUserAndPost(t, app, userHandler, params)

	authHandler := NewAuthHandler(userStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "foo@bar.com",
		Password: "foobar123",
	}

	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status of 200 but got %d ", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}
	if authResp.Token == "" {
		t.Fatal("expected the JWT token in the auth response")
	}
}

func TestAuthWithWrongPassword(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

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
	_ = createUserAndPost(t, app, userHandler, params)

	authHandler := NewAuthHandler(userStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "foo@bar.com",
		Password: "foobaz123",
	}

	b, _ := json.Marshal(authParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == 200 {
		t.Fatal("not expected status code of 200")
	}
}
