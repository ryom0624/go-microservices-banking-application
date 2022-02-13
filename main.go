package main

import (
	"go-microservices-banking-application/app"
	"go-microservices-banking-application/logger"
)

func main() {

	logger.Info("Starting the application")

	app.Start()

}
