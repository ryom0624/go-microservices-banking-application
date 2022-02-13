package main

import (
	"banking/app"
	"local.packages/logger"
)

func main() {

	logger.Info("Starting the application")

	app.Start()

}
