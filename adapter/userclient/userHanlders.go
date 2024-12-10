package userclient

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"scarlet_backend/adapter"
	"scarlet_backend/model"
)

var userService adapter.UserService

func InitUserService(service adapter.UserService) { userService = service }

func handlePanic(resp http.ResponseWriter) {
	if r := recover(); r != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Internal Server Error}`))
	}
}

func GetUsers(resp http.ResponseWriter, req *http.Request) {
	defer handlePanic(resp)
	resp.Header().Set("Content-type", "application/json")
	users, err := userService.FindAll()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error getting users"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(users)
}

func GetUserById(resp http.ResponseWriter, req *http.Request) {
	defer handlePanic(resp)
	resp.Header().Set("Content-type", "application/json")
	var params map[string]int
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}
	id := params["id"]
	user, err := userService.FindById(id)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error getting user"}`))
		return
	}
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error:" "Usuario not found"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func GetUserByEmail(resp http.ResponseWriter, req *http.Request) {
	defer handlePanic(resp)
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	email := params["email"]
	user, err := userService.FindByEmail(email)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error getting the user"}`))
		return
	}
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error": "User not found"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func GetUserByPhone(resp http.ResponseWriter, req *http.Request) {
	defer handlePanic(resp)
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	phone := params["phone"]
	user, err := userService.FindByPhone(phone)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error getting user"}`))
		return
	}
	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"error:" "Usuario not found"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}

func SaveByEmail(resp http.ResponseWriter, req *http.Request) {
	defer handlePanic(resp)
	resp.Header().Set("Content-type", "application/json")

	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Invalid request payload"}`))
		return
	}

	existingUser, err := userService.FindByEmail(user.Email)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error checking user"}`))
		return
	}
	if existingUser != nil {
		resp.WriteHeader(http.StatusConflict)
		resp.Write([]byte(`{"error": "Email already exists"}`))
		return
	}

	user.Id = rand.Int63()
	savedUser, err := userService.SaveByEmail(&user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error saving user"}`))
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(savedUser)
}

func CheckLogin(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	email := params["email"]
	password := params["password"]
	userFinded, err := userService.CheckLogin(email, password)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error verificando el usuario"}`))
		return
	}
	if userFinded == nil {
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

func SendOTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	phone := params["phone"]
	otp_id, err := userService.SendOTP(phone)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error enviando el OTP"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]string{"otp_id": otp_id})
}

func VerifyOTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")
	var params map[string]string
	json.NewDecoder(req.Body).Decode(&params)
	otp_id := params["otp_id"]
	otp_code := params["otp_code"]
	err := userService.VerifyOTP(otp_id, otp_code)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error": "Error verificando el OTP: ` + err.Error() + `", "otp_id": "` + otp_id + `", "otp_code": "` + otp_code + `"}`))
		return
	}
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]string{"status": "verified"})
}

func AddUsersByPhone(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")
	var user model.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{"error:" "Error marshalling the request"}`))
		return
	}
	user.Id = rand.Int63()
	userService.SaveByPhone(&user)
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(user)
}
