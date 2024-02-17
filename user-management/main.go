package main

import (
	"flag"
	"fmt"

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
	app.Listen(listenAddr)

}
