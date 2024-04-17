package main

import (
	"scarlet_backend/internal/domain/entities"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"log"
)

func main(){
	router := mux.NewRouter();
	const port string = ":8000"
	router.HandleFunc("/", func(resp  http.ResponseWriter, req *http.Request){
		fmt.Fprintln(resp, "Up and running...")
	})
	router.HandleFunc("/users", entities.GetUsers).Methods("GET")
	router.HandleFunc("/users", entities.AddUsers).Methods("POST")
	log.Println("Server listening on port ", port)
	log.Fatalln(http.ListenAndServe(port, router))
}