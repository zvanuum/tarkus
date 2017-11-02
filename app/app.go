package app

import (
	"net/http"
	"fmt"
	"log"
	
	"github.com/gorilla/mux"
)

type Server struct {
	Port int
}

func (app *Server) ServeHTTP() {
	r := mux.NewRouter()
	
	app.InitializeRoutes(r)
	
	log.Printf("Server listening on port %d", app.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", app.Port), r)
}

func (app *Server) InitializeRoutes(r *mux.Router) {
	r.HandleFunc("/chain", test).Methods("GET")
	r.HandleFunc("/mine", test).Methods("GET")

	r.HandleFunc("/transaction/new", test).Methods("POST")
}

func test(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
}
