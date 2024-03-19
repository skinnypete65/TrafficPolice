package main

import (
	_ "TrafficPolice/docs"
	"TrafficPolice/internal/app"
)

// @title Traffic Police API
// @version 1.0
// @description API Server for Traffic Police Application

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
