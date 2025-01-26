package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	sandbox.SandboxPtr = sandbox.NewSandbox(10)
	defer sandbox.SandboxPtr.Cleanup()

	// 設置信號處理
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Received interrupt signal, cleaning up...")
		sandbox.SandboxPtr.Cleanup()
		os.Exit(0)
	}()

	sandbox.SandboxPtr.RunShellCommandByRepo("user_name/repo_name", nil)

	database.DBConn.AutoMigrate(&models.Book{})
	database.DBConn.AutoMigrate(&models.Sandbox{})

	app := routes.New()
	log.Fatal(app.Listen(":3001"))
}
