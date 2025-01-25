package main

import (
	"log"

	"github.com/HeavenManySugar/OJ-PoC/database"
	"github.com/HeavenManySugar/OJ-PoC/models"
	"github.com/HeavenManySugar/OJ-PoC/routes"
	"github.com/HeavenManySugar/OJ-PoC/sandbox"
)

// @title			OJ-PoC API
// @version		1.0
// @description	This is a simple OJ-PoC API server.
// @BasePath		/
func main() {
	if err := database.Connect(); err != nil {
		log.Panic("Can't connect database:", err.Error())
	}
	s := sandbox.NewSandbox(10)
	s.RunShellCommand([]byte("echo 'Hello, World!'"))

	database.DBConn.AutoMigrate(&models.Book{})

	app := routes.New()
	log.Fatal(app.Listen(":3001"))
}
