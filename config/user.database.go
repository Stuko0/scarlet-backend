package repository

import (
	"context"
	"fmt"
	"log"
	"scarlet_backend/internal/domain/entities"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type UserRepository interface{
	SaveByEmail(user *entities.User)(*entities.User, error)
	FindAll()([]entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	CheckLogin(email string, psw string) (*entities.User, error)
	FindByPhone(phone string) (*entities.User, error)
	SendVerificationCode(phone string, code string) error
}

type repo struct{}

func NewUserRepository() UserRepository{
	return &repo{}
}

const (
	projectId 	string = "scarlet-419401"
	collectionName string = "users"
)

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
		"psw": user.Pwd,
		"origin": user.Origin,
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
			Pwd: doc.Data()["psw"].(string),
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

func (r *repo) CheckLogin(email string, psw string) (*entities.User, error){
	user, err := r.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("no user found with email %s", email)
	}
	if user.Pwd != psw {
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

func (*repo) SendVerificationCode(phone string, code string) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile("D:/clases/Integrador/backend/scarlet-419401.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error al inicializar la aplicación de Firebase: %v", err)
		return err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("Error al obtener el cliente de Messaging: %v", err)
		return err
	}
	// Crear el mensaje
	message := &messaging.Message{
		Data: map[string]string{
			"code": code,
		},
		Token: phone,
	}

	// Enviar el mensaje
	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("Error al enviar el mensaje: %v", err)
		return err
	}

	// Imprimir el ID del mensaje
	log.Printf("Mensaje enviado con éxito, ID: %s", response)
	return nil
}
