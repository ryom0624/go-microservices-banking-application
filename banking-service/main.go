package main

import (
	"banking/app"
	"local.packages/lib/logger"
)

func main() {

	logger.Info("Starting the application")

	app.Start()

}
