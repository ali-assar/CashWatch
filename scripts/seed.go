package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ali-Assar/CashWatch/authentication-service/db"
	"github.com/Ali-Assar/CashWatch/authentication-service/handler"
	"github.com/Ali-Assar/CashWatch/types"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	_, err = database.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Fatal(err)
	}
	db.CreateTable(database)

	userStore := db.NewPostgreSQLUserStore(database)

	user1 := addUser(userStore, "james", "foo")
	fmt.Println("user token ->", handler.CreateTokenFromUser(user1))
	user2 := addUser(userStore, "bar", "baz")
	fmt.Println("user token ->", handler.CreateTokenFromUser(user2))
}

func addUser(store db.UserStorer, fn, ln string) *types.User {
	user, err := types.NewUserFormParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
