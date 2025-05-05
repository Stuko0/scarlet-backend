package wildfire

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WildfireNRTRepository interface{
	UpsertFromDetection(ctx context.Context, detection *WildfireNRT)error
	GetActiveFires(ctx context.Context)([]*WildfireNRT, error)
	GetNearbyFires(ctx context.Context, lat, lng float64, radiusKm int)([]*WildfireNRT, error)
	BulkUpsertFromDetections(ctx context.Context, detections[]*WildfireNRT)error
}

type MongoRepository struct{
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database)*MongoRepository{
	return &MongoRepository{
		collection: db.Collection("current_wildfires"),
	}
}

func (r *MongoRepository) UpsertFromDetection(ctx context.Context, detection *WildfireNRT) error{
	update:=bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"status": "active",
			"weather": detection.Weather,
			"metadata": detection.MetaData,
		},
		"$setOnInsert": bson.M{
			"fire_id": detection.FireID,
			"source": detection.Source,
			"external_id": detection.ExternalID,
			"discovery_time": detection.DiscoveredTime,
			"location": detection.Location,
		},
	}
	_, err:= r.collection.UpdateOne(
		ctx, bson.M{"fire_id":detection.FireID},update, options.Update().SetUpsert(true),
	)
	return err
}

func(r *MongoRepository)BulkUpsertFromDetections(ctx context.Context, detections []*WildfireNRT)error{
	if len(detections)==0{return nil}

	models:= make([]mongo.WriteModel, len(detections))
	for i, dedetection:=range detections{
		update:=bson.M{
			"$set": bson.M{
				"last_updated": dedetection.LastUpdated,
				"status": "active",
				"weather": dedetection.Weather,
				"metadata": dedetection.MetaData,
			},
			"$setOnInsert": bson.M{
				"fire_id":dedetection.FireID,
				"source":dedetection.Source,
				"external_id":dedetection.ExternalID,
				"discovery_time":dedetection.DiscoveredTime,
				"location": dedetection.Location,
			},
		}
		models[i]=mongo.NewUpdateOneModel().
			SetFilter(bson.M{"fire_id":dedetection.FireID}).
			SetUpdate(update).
			SetUpsert(true)
	}

	_, err:=r.collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return err
}

func (r *MongoRepository) GetNearbyFires (ctx context.Context, lat, lng float64, radiusKm int)([]*WildfireNRT, error){
	filter:=bson.M{
		"location":bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type": "Point",
					"coordinates": []float64{lat,lng},
				},
				"$maxDistance": radiusKm*1000,
			},
		},
	}
	cursor, err:=r.collection.Find(ctx, filter)
	if err!=nil{
		return nil, fmt.Errorf("failed to find nearby fires: %w", err)
	}
	var fires []*WildfireNRT
	if err := cursor.All(ctx, &fires); err!=nil{
		return nil, fmt.Errorf("failed to decode nearby fires: %w", err)
	}
	return fires, nil
}

func (r *MongoRepository) GetActiveFires(ctx context.Context)([]*WildfireNRT, error){
	filter:=bson.M{
		"status": "active",
	}
	cursor, err:=r.collection.Find(ctx, filter)
	if err!=nil{
		return nil, fmt.Errorf("failed to find active fires: %w", err)
	}
	var fires[]*WildfireNRT
	if err:= cursor.All(ctx, &fires); err!=nil{
		return nil, fmt.Errorf("failed to decode active fires: %w", err)
	}
	return fires, nil
}