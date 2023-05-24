package main

import (
	"github.com/joho/godotenv"
	"proxy-harvester/internal/pkg/manager"
	"proxy-harvester/internal/server"
)

func main() {
	if err := godotenv.Load("settings.env"); err != nil {
		panic("Error loading .env file")
	}

	go manager.Start()
	server.Start()
}
