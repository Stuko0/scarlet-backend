package services

import (
	"encoding/json"
	"math/rand"
	"net/http"
	repository "scarlet_backend/config"
	"scarlet_backend/internal/domain/entities"
	"strconv"
	"time"
)

var(repo repository.UserRepository = repository.NewUserRepository())

func GetUsers(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	users, err :=  repo.FindAll()
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo usuarios"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(users)
}

func GetUserByEmail(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	email := params["email"]
	user, err :=  repo.FindByEmail(email)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo el usuario"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func GetUserByPhone(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	phone := params["phone"]
	user, err :=  repo.FindByPhone(phone)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo el usuario"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func AddUsers(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var user entities.User
	err:=json.NewDecoder(req.Body).Decode(&user)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error marshalling the request"}`))
		return
	}
	user.Id=rand.Int63()
	user.Origin="email"
	repo.SaveByEmail(&user)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func CheckLogin(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	email := params["email"]
	password := params["password"]
	userFinded, err :=  repo.CheckLogin(email, password)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error verificando el usuario"}`))
		return
	}
	if userFinded  ==  nil {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte(`{"error:" "Correo electrónico o contraseña incorrectos"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"status:" "Inicio de sesión exitoso"}`))
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // Genera un número aleatorio de 6 dígitos
	return strconv.Itoa(code)
}

func SendVerificationCode(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	phone := params["phone"]
	code := params[generateVerificationCode()]

	err :=  repo.SendVerificationCode(phone, code)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error enviando el código de verificación"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(entities.User{Phone: phone})
}