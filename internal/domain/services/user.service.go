package services

import (
	"encoding/json"
	"math/rand"
	"net/http"
	repository "scarlet_backend/config"
	"scarlet_backend/internal/domain/entities"
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

func GetUserById(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]int
	json.NewDecoder(req.Body).Decode(&params)
	id := params["id"]
	user, err :=  repo.FindById(id)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error obteniendo el usuario"}`))
		return
	}
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error:" "Usuario no encontrado"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
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
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error:" "Usuario no encontrado"}`))
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
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error:" "Usuario no encontrado"}`))
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
	userJson, err := json.Marshal(userFinded)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error serializando el usuario"}`))
		return
	}
	resp.Write(userJson)
}

func SendOTP(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	phone := params["phone"]
	otp_id, err :=  repo.SendOTP(phone)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error enviando el OTP"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]string{"otp_id": otp_id})
}

func VerifyOTP(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	otp_id := params["otp_id"]
	otp_code := params["otp_code"]
	err :=  repo.VerifyOTP(otp_id, otp_code)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error verificando el OTP: ` + err.Error() + `", "otp_id": "` + otp_id + `", "otp_code": "` + otp_code + `"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]string{"status": "verified"})
}

func AddUsersByPhone(resp http.ResponseWriter, req *http.Request){
	resp.Header().Set("Content-type", "application/json")
	var user entities.User
	err:=json.NewDecoder(req.Body).Decode(&user)
	if err!= nil{
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error marshalling the request"}`))
		return
	}
	user.Id=rand.Int63()
	repo.SaveByPhone(&user)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}
