package main

import (
	"flag"
	
	"github.com/zachvanuum/tarkus/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port for server to listen on")
	flag.Parse()
	
	app := app.Server{ Port: port}
	app.ServeHTTP()
}
