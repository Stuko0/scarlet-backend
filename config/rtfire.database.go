package repository

import (
	"context"
	"log"
	"scarlet_backend/internal/domain/entities"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type RTFireRepository interface{
	SaveRTFire(fire *entities.RTFire)(*entities.RTFire, error)
	GetRTFireFromDB() ([]*entities.RTFire, error)
	DeleteAllRTFires() error
}


type rtfireRepo struct{}

func NewRTFireRepository() RTFireRepository{
	return &rtfireRepo{}
}

func (*rtfireRepo) SaveRTFire(fire *entities.RTFire) (*entities.RTFire, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	_,_, err=client.Collection("rtfires").Add(ctx, map[string]interface{}{
		"id":fire.Id,
		"latitude": fire.Latitude,
		"longitude": fire.Longitude,
		"detectedAt": fire.DetectedAt,
		"confidence": fire.Confidence,
		"frp": fire.FRP,
		"fwi": fire.FWI,
		"fireType": fire.FireType,
	})

	if err != nil{
		log.Fatalf("No se pudo registrar el incendio: %v", err)
		return nil, err
	}
	return fire, nil
}

func (*rtfireRepo) GetRTFireFromDB() ([]*entities.RTFire, error){
	ctx := context.Background()
	client, err  :=  firestore.NewClient(ctx, projectId)
	if err != nil{
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return nil, err
	}
	
	defer client.Close()
	
	var fires []*entities.RTFire
	iter := client.Collection("rtfires").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error obteniendo los incendios: %v", err)
			return nil, err
		}
		var fire entities.RTFire
		doc.DataTo(&fire)
		fires = append(fires, &fire)
	}

	return fires, nil
}

func (*rtfireRepo) DeleteAllRTFires() error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Printf("No se pudo crear la conexion a la base de datos: %v", err)
		return err
	}
	defer client.Close()

	iter := client.Collection("rtfires").Documents(ctx)
	batch := client.Batch()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error obteniendo los incendios: %v", err)
			return err
		}
		batch.Delete(doc.Ref)
	}

	_, err = batch.Commit(ctx)
	if err != nil {
		log.Printf("Error eliminando los incendios: %v", err)
		return err
	}

	return nil
}
