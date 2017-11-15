package app

import (
	"net/http"
	"fmt"
	"log"
	
	"github.com/gorilla/mux"
	"github.com/zachvanuum/tarkus/blockchain"
	"github.com/zachvanuum/tarkus/api"
)

type App struct {
	Blockchain *blockchain.Blockchain
}

func (app *App) ServeHTTP(port int) {
	r := mux.NewRouter()
	
	InitializeRoutes(app, r)
	
	log.Printf("Server listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func InitializeRoutes(app *App, r *mux.Router) {
	r.HandleFunc("/chain", api.ChainHandler(app.Blockchain)).Methods("GET")
	r.HandleFunc("/mine", api.MineHandler(app.Blockchain)).Methods("GET")
	r.HandleFunc("/node/consensus", api.ConsensusHandler(app.Blockchain)).Methods("GET")

	r.HandleFunc("/node/register", api.RegisterNodesHandler(app.Blockchain)).Methods("POST")
	r.HandleFunc("/transaction/new", api.NewTransactionHandler(app.Blockchain)).Methods("POST")
}
