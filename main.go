package main

import (
	"log"
	"net/http"
	"os"
	"scarlet_backend/internal/domain/services"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err!= nil{log.Fatalf("No se pudo encontrar al usuario: %v", err)}
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if credentialsFile == "" {
        log.Printf("ADVERTENCIA: GOOGLE_APPLICATION_CREDENTIALS no está configurado")
        // Configurar explícitamente si no está en el entorno
        os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./scarlet-419401.json")
    }
	log.Printf("GOOGLE_APPLICATION_CREDENTIALS: %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	ticker := time.NewTicker(60 * time.Minute)
	quit := make(chan struct{})
	router := mux.NewRouter()
	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("Up and running..."))
	})
	router.HandleFunc("/users", services.GetUsers).Methods("GET")
	router.HandleFunc("/users/register", services.AddUsers).Methods("POST")
	router.HandleFunc("/users/getUserByEmail", services.GetUserByEmail).Methods("POST")
	router.HandleFunc("/users/getUserById", services.GetUserById).Methods("POST")
	router.HandleFunc("/users/checkLogin", services.CheckLogin).Methods("POST")
	router.HandleFunc("/users/sendPhoneCode", services.SendOTP).Methods("POST")
	router.HandleFunc("/users/verifyPhoneCode", services.VerifyOTP).Methods("POST")
	router.HandleFunc("/users/getUserByPhone", services.GetUserByPhone).Methods("POST")
	router.HandleFunc("/users/registerUserByPhoneNumber", services.AddUsersByPhone).Methods("POST")
	router.HandleFunc("/fires", services.GetSavedFires).Methods("GET")
	router.HandleFunc("/fires/save", services.SaveFire).Methods("POST")
	router.HandleFunc("/rtfires", services.GetSavedRTFires).Methods("GET")
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := services.SaveRTFireData()
				if err != nil {
					log.Printf("Error guardando los datos del incendio: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	router.HandleFunc("/rtfires/save", services.SaveRTFire).Methods("POST")
	router.HandleFunc("/rtfires/delete", services.DeleteAllRTFires).Methods("DELETE")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	const port string = ":8000"
	log.Println("Server listening on port ", port)

	// Inicia el servidor con el middleware de CORS
	log.Fatalln(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
