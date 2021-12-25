package main

import (
	_config "github.com/fanajib5/url_shortener_PA/config"
	_routes "github.com/fanajib5/url_shortener_PA/routes"
)

func main() {
	_config.InitDB()
	e := _routes.New()

	// Run REST API Server
	// e.Start(":8080")
	e.Logger.Fatal(e.Start(":8080"))
}
