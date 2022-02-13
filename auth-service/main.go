package main

import (
	"auth/app"
	"local.packages/lib/logger"
)

func main() {

	logger.Info("Starting banking authentication application")

	app.Start()

}
