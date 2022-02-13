package main

import (
	"banking/app"
	"local.packages/lib/logger"
)

func main() {

	logger.Info("Starting banking application")

	app.Start()

}
