package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jmcerez0/fiber-demo/handlers"
	"github.com/jmcerez0/fiber-demo/middlewares"
	"github.com/jmcerez0/fiber-demo/utils"
)

func init() {
	utils.LoadEnv()
	utils.CreateDB()
	utils.ConnectToDB()
	utils.MigrateSchema()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/signup", handlers.SignUp)

	app.Post("/signin", handlers.SignIn)

	app.Get("/users", middlewares.RequireAuth, handlers.GetAllUsers)

	app.Listen(":" + os.Getenv("PORT"))
}
