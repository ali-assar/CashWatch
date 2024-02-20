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
	defer tdb.Exec("DROP TABLE users")
	// defer tdb.Close()
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

func TestGetUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

	// defer tdb.Close()
	userStore := db.NewPostgreSQLUserStore(tdb)
	app := fiber.New()

	userHandler := NewUserHandler(userStore)
	app.Post("/", userHandler.HandlePostUser)

	insertedUser := types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foo@bar.com",
		Password:  "foobar123",
	}
	user := createUserAndPost(t, app, userHandler, insertedUser)

	app.Get("/:id", userHandler.HandleGetUserByID)

	req := httptest.NewRequest("GET", "/"+"1", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}
	var retrievedUser types.User
	err = json.NewDecoder(resp.Body).Decode(&retrievedUser)
	if err != nil {
		t.Fatal(err)
	}
	if retrievedUser.ID != user.ID {
		t.Errorf("expected the retrieved user ID to be %d, but got %d", user.ID, retrievedUser.ID)
	}
	if retrievedUser.FirstName != user.FirstName {
		t.Errorf("expected the retrieved user's first name to be %s, but got %s", user.FirstName, retrievedUser.FirstName)
	}
	if retrievedUser.LastName != user.LastName {
		t.Errorf("expected the retrieved user's last name to be %s, but got %s", user.LastName, retrievedUser.LastName)
	}
}

func TestDeleteUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

	// defer tdb.Close()
	userStore := db.NewPostgreSQLUserStore(tdb)
	app := fiber.New()

	userHandler := NewUserHandler(userStore)
	app.Post("/", userHandler.HandlePostUser)

	insertedUser := types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foo@bar.com",
		Password:  "foobar123",
	}
	createUserAndPost(t, app, userHandler, insertedUser)

	app.Delete("/:id", userHandler.HandleGetUserByID)
	req := httptest.NewRequest("DELETE", "/"+"1", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}
}

func TestUpdateUser(t *testing.T) {
	tdb, err := db.InitDB()
	if err != nil {
		t.Errorf("database error")
	}
	defer tdb.Exec("DROP TABLE users")

	// defer tdb.Close()
	userStore := db.NewPostgreSQLUserStore(tdb)
	app := fiber.New()

	userHandler := NewUserHandler(userStore)
	app.Post("/", userHandler.HandlePostUser)
	insertedUser := types.CreateUserParams{
		FirstName: "foo",
		LastName:  "bar",
		Email:     "foo@bar.com",
		Password:  "foobar123",
	}
	user := createUserAndPost(t, app, userHandler, insertedUser)

	app.Put("/:id", userHandler.HandleGetUserByID)
	updatedParams := &types.CreateUserParams{
		FirstName: "poo",
		LastName:  "baz",
	}
	b, _ := json.Marshal(updatedParams)

	req := httptest.NewRequest("PUT", "/1", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}

	app.Get("/:id", userHandler.HandleGetUserByID)
	req = httptest.NewRequest("GET", "/"+"1", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var updatedUser types.User
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	if err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}
	if updatedUser.FirstName == updatedParams.FirstName {
		t.Errorf("expected the first name to be updated, but it's still %s", user.FirstName)
	}
	if updatedUser.LastName == updatedParams.LastName {
		t.Errorf("expected the last name to be updated, but it's still %s", user.LastName)
	}
}
