package app

import (
	"net/http"
	"fmt"
	"log"
	
	"github.com/gorilla/mux"
	"github.com/zachvanuum/tarkus/blockchain"
)

type App struct {
	Blockchain blockchain.Blockchain
}

func (app *App) ServeHTTP(port int) {
	r := mux.NewRouter()
	
	InitializeRoutes(r)
	
	log.Printf("Server listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func InitializeRoutes(r *mux.Router) {
	r.HandleFunc("/chain", test).Methods("GET")
	r.HandleFunc("/mine", test).Methods("GET")

	r.HandleFunc("/transaction/new", test).Methods("POST")
}

func test(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
}
