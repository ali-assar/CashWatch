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

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	listenAddr := flag.String("HTTP listenAddr", ":3000", "the listen address of HTTP server")
	flag.Parse()

	fmt.Println("server is listening at", *listenAddr)

	var (
		userStore   = db.NewPostgreSQLUserStore(database)
		userHandler = handler.NewUserHandler(userStore)
		authHandler = handler.NewAuthHandler(userStore)

		app         = fiber.New()
		apiRegister = app.Group("/")
		auth        = app.Group("/api")
		apiv1       = app.Group("/api/v1", handler.JWTAuthentication(userStore))
	)

	//auth
	auth.Post("/auth", authHandler.HandleAuthenticate)
	//user
	apiRegister.Post("/register", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:email", userHandler.HandleDeleteUser)

	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
