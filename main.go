package main

import (
	"github.com/fanajib5/url_shortener_PA/config"
	"github.com/fanajib5/url_shortener_PA/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	// Run REST API Server
	// e.Start(":8080")
	e.Logger.Fatal(e.Start(":8080"))
}
