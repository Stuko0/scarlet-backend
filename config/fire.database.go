package repository

import (
	"context"
	"log"
	"scarlet_backend/internal/domain/entities"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FireRepository interface{
	SaveFire(fire *entities.Fire)(*entities.Fire, error)
	GetFireFromDB() ([]*entities.Fire, error)
}

type fireRepo struct{}

func NewFireRepository() FireRepository{
	return &fireRepo{}
}

func (*fireRepo) SaveFire(fire *entities.Fire) (*entities.Fire, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	_,_, err=client.Collection("fires").Add(ctx, map[string]interface{}{
		"id":fire.Id,
		"latitude": fire.Latitude,
		"longitude": fire.Longitude,
		"acq_date": fire.AcqDate,
		"acq_time": fire.AcqTime,
		"daynight": fire.DayNight,
	})

	if err != nil{
		log.Fatalf("No se pudo registrar el incendio: %v", err)
		return nil, err
	}
	return fire, nil
}

func (*fireRepo) GetFireFromDB() ([]*entities.Fire, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	
	var fires []*entities.Fire
	iter := client.Collection("fires").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error obteniendo los incendios: %v", err)
			return nil, err
		}
		var fire entities.Fire
		doc.DataTo(&fire)
		fires = append(fires, &fire)
	}

	return fires, nil
}