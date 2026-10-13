package main

import (
	"banking-auth/app"
	"banking-auth/logger"
)

func main() {
	// log.Println("Starting App...")
	logger.Info("Starting App")
	app.Start()
}
