package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Ali-Assar/CashWatch/user-management/db"
	"github.com/Ali-Assar/CashWatch/user-management/handler"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// TODO FIX this main
func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	userStore := db.NewPostgreSQLUserStore(database)

	userHandler := handler.NewUserHandler(userStore)
	authHandler := handler.NewAuthHandler(userStore)

	listenAddr := flag.String("HTTP listenAddr", ":3000", "the listen address of HTTP server")
	flag.Parse()

	fmt.Println("server is listening at", *listenAddr)

	app := fiber.New()
	auth := app.Group("/login")
	apiRegister := app.Group("/")
	apiv1 := app.Group("/api/v1", handler.JWTAuthentication(userStore))

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiRegister.Post("/register", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:email", userHandler.HandleDeleteUser)
	//TODO: implement get user by email

	// Start Fiber app
	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}

}
