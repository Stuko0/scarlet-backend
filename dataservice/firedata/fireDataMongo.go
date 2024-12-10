package firedata

import (
	"context"
	"scarlet_backend/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FireDataMongo struct {
	db *mongo.Database
}

func NewFireDataMongo(client *mongo.Client, dbName string) *FireDataMongo {
	return &FireDataMongo{
		db: client.Database(dbName),
	}
}

func (f *FireDataMongo) GetFires() ([]model.RTFire, error) {
	var fires []model.RTFire
	collection := f.db.Collection("rtfires")
	cursor, err := collection.Find(context.Background(), options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var fire model.RTFire
		if err = cursor.Decode(&fire); err != nil {
			return nil, err
		}
		fires = append(fires, fire)
	}
	return fires, nil
}

func (f *FireDataMongo) AddFire(fire model.RTFire) error {
	collection := f.db.Collection("rtfires")
	_, err := collection.InsertOne(context.Background(), fire)
	return err
}
