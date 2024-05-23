package main

import (
	"scarlet_backend/internal/domain/services"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
)

func main(){
	router := mux.NewRouter();
	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request){
		resp.Write([]byte("Up and running..."))
	})
	router.HandleFunc("/users", services.GetUsers).Methods("GET")
	router.HandleFunc("/users/register", services.AddUsers).Methods("POST")
	router.HandleFunc("/users/getUserByEmail", services.GetUserByEmail).Methods("POST")
	router.HandleFunc("/users/checkLogin", services.CheckLogin).Methods("POST")
	router.HandleFunc("/users/sendPhoneCode", services.SendOTP).Methods("POST")
	router.HandleFunc("/users/verifyPhoneCode", services.VerifyOTP).Methods("POST")
	router.HandleFunc("/users/getUserByPhone", services.GetUserByPhone).Methods("POST")
	router.HandleFunc("/users/registerUserByPhoneNumber", services.AddUsersByPhone).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	const port string = ":8000"
	log.Println("Server listening on port ", port)

	// Inicia el servidor con el middleware de CORS
	log.Fatalln(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}