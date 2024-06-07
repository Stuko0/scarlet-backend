package repository

import (
	"context"
	"fmt"
	"log"
	"scarlet_backend/internal/domain/entities"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type UserRepository interface{
	SaveByEmail(user *entities.User)(*entities.User, error)
	FindAll()([]entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindById(id int) (*entities.User, error)
	CheckLogin(email string, psw string) (*entities.User, error)
	SendOTP(phone string)(string, error)
	VerifyOTP(otp_id string, otp_code string) error
	FindByPhone(phone string) (*entities.User, error)
	SaveByPhone(user *entities.User) (*entities.User, error)
}

type repo struct{}

func NewUserRepository() UserRepository{
	return &repo{}
}

const (
	projectId 	string = "scarlet-419401"
	collectionName string = "users"
)

type OTPResponse struct {
	OtpID string `json:"otp_id"`
}

func (*repo) SaveByEmail(user *entities.User) (*entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	_,_, err=client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"id": user.Id,
		"name": user.Name,
		"lastname": user.Lastname,
		"email": user.Email,
		"psw": user.Psw,
		"origin": "email",
	})

	if err != nil{
		log.Fatalf("No se pudo registrar al usuario: %v", err)
		return nil, err
	}
	return user, nil
}

func (*repo) FindAll()([]entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Fatalf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	defer client.Close()
	var users[]entities.User
	itr:=client.Collection(collectionName).Documents(ctx)
	for{
		doc, err := itr.Next()
		if err == iterator.Done{break}
		if err!=nil{
			log.Fatalf("No se pudo cargar la informacion de los usuarios: %v", err)
			return nil, err
		}
		user := entities.User{
			Id:	doc.Data()["id"].(int64),
			Name: doc.Data()["name"].(string),
			Lastname: doc.Data()["lastname"].(string),
			Email: doc.Data()["email"].(string),
			Psw: doc.Data()["psw"].(string),
			Origin: doc.Data()["origin"].(string),
		}
		users = append(users, user)
	}
	return users, nil
}

func (*repo) FindByEmail(email string) (*entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	query := client.Collection(collectionName).Where("email", "==", email).Where("origin", "==", "email")
	dsnap, err := query.Documents(ctx).Next()
	if err == iterator.Done {
		log.Println("No se encontró ningún usuario con ese correo electrónico")
		return nil, fmt.Errorf("No se encontró ningún usuario con ese correo electrónico")
	} else
	if err != nil {
		log.Fatalf("No se pudo encontrar al usuario: %v", err)
		return nil, err
	}
	var user entities.User
	if err := dsnap.DataTo(&user); err != nil {
    log.Fatalf("Error mapping document to User: %v", err)
    return nil, err
}
	dsnap.DataTo(&user)
	return &user, nil
}

func (*repo) FindById(id int) (*entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	query := client.Collection(collectionName).Where("id", "==", id)
	dsnap, err := query.Documents(ctx).Next()
	if err == iterator.Done {
		log.Println("No se encontró ningún usuario con ese id")
		return nil, fmt.Errorf("No se encontró ningún usuario con ese id")
	} else
	if err != nil {
		log.Fatalf("No se pudo encontrar al usuario: %v", err)
		return nil, err
	}
	var user entities.User
	if err := dsnap.DataTo(&user); err != nil {
    log.Fatalf("Error mapping document to User: %v", err)
    return nil, err
}
	dsnap.DataTo(&user)
	return &user, nil
}

func (r *repo) CheckLogin(email string, psw string) (*entities.User, error){
	user, err := r.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("no user found with email %s", email)
	}
	if user.Psw == "" {
		return nil, fmt.Errorf("password not set for email %s", email)
	}
	if user.Psw != psw {
		return nil, fmt.Errorf("incorrect password for email %s", email)
	}
	return user, nil
}

func (*repo) FindByPhone(phone string) (*entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	query := client.Collection(collectionName).Where("phone", "==", phone).Where("origin", "==", "phone")
	dsnap, err := query.Documents(ctx).Next()
	if err == iterator.Done {
		log.Println("No se encontró ningún usuario con ese número de teléfono")
		return nil, nil
	} else
	if err != nil {
		log.Fatalf("No se pudo encontrar al usuario: %v", err)
		return nil, err
	}
	var user entities.User
	dsnap.DataTo(&user)
	return &user, nil
}

func (*repo) SendOTP(phone string) (string, error) {
	url := "https://d7sms.p.rapidapi.com/verify/v1/otp/send-otp"
	payload := map[string]string{
		"originator": "SignOTP",
		"recipient": phone,
		"content": "Tu codigo de Scarlet App es: {}",
		"expiry": "600",
		"data_coding": "text",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoLWJhY2tlbmQ6YXBwIiwic3ViIjoiYTViMDAwYWEtODZmOS00YzcyLThhNzItMTI0N2NkN2E3MzdkIn0.evTIuqF8ixasgJF39M1aS9ohoUGX2wHi9yN80lQkQLc")
	req.Header.Add("X-RapidAPI-Key", "a7b5094c54msh978fed74d9e9925p1bd538jsncaf858139020")
	req.Header.Add("X-RapidAPI-Host", "d7sms.p.rapidapi.com")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var otpResponse OTPResponse
	json.Unmarshal(body, &otpResponse)
	return otpResponse.OtpID, nil
}

func (*repo) VerifyOTP(otp_id string, otp_code string) error {
	url := "https://d7sms.p.rapidapi.com/verify/v1/otp/verify-otp"
	payload := map[string]string{
		"otp_id": otp_id,
		"otp_code": otp_code,
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoLWJhY2tlbmQ6YXBwIiwic3ViIjoiYTViMDAwYWEtODZmOS00YzcyLThhNzItMTI0N2NkN2E3MzdkIn0.evTIuqF8ixasgJF39M1aS9ohoUGX2wHi9yN80lQkQLc")
	req.Header.Add("X-RapidAPI-Key", "a7b5094c54msh978fed74d9e9925p1bd538jsncaf858139020")
	req.Header.Add("X-RapidAPI-Host", "d7sms.p.rapidapi.com")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Error al verificar el OTP: %s", res.Status)
	}
	return nil
}

func (*repo) SaveByPhone(user *entities.User) (*entities.User, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	_,_, err=client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"id": user.Id,
		"name": user.Name,
		"lastname": user.Lastname,
		"email": user.Email,
		"phone": user.Phone,
		"origin": "phone",
	})

	if err != nil{
		log.Fatalf("No se pudo registrar al usuario: %v", err)
		return nil, err
	}
	return user, nil
}
