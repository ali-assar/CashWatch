package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber"
)

func main() {
	listenAddr := flag.String("HTTP listenAddr", ":3000", "the listen address of HTTP server")
	flag.Parse()

	fmt.Println("server is listening at", ":3000")

	app := fiber.New()

	app.Post("/api/user/register", RegisterHandler)
	app.Post("/api/user/login", LoginHandler)
	app.Get("/api/user/profile", ProfileHandler)

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

	// Start Fiber app
	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}

}
