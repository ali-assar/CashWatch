package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber"
)

// TODO FIX this main
func main() {
	// Establish database connection
	db, err := sql.Open("postgres", "postgres://admin:admin@localhost:5432/user")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userStore = db.NewMongoUserStore(client)
	store = &db.Store{User: userStore}
	userHandler = api.NewUserHandler(userStore)

	listenAddr := flag.String("HTTP listenAddr", ":3000", "the listen address of HTTP server")
	flag.Parse()
	fmt.Println("server is listening at", ":3000")

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user/", UserHandler.HandlePostUser)

	// Start Fiber app
	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}

}
