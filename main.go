package main

import (
	"flag"
	
	"github.com/zachvanuum/tarkus/app"
	"github.com/zachvanuum/tarkus/blockchain"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port for server to listen on")
	flag.Parse()
	
	blockchain := blockchain.InitializeBlockchain()
	
	app := app.App{Blockchain: blockchain}
	app.ServeHTTP(port)
}
