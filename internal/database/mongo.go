package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient (ctx context.Context)(*mongo.Client, error){
	client, err:= mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetServerSelectionTimeout(5 *time.Second).
		SetConnectTimeout(10 *time.Second).
		SetSocketTimeout(15 *time.Second))
	
	if err!= nil{
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err:= client.Ping(ctx, readpref.Primary()); err!=nil{
		return nil, fmt.Errorf("failed to ping MongoDB: %w",err)
	}
	return client, nil
}

func EnsureMongoCollections(ctx context.Context, client *mongo.Client)error{
	wildfireDb:=client.Database("wildfire")

	_,err:= wildfireDb.Collection("current").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:bson.D{{Key: "location", Value: "2dsphere"}},
	})

	if err!= nil{return fmt.Errorf("failed to create geo index: %w",err)}
	_, err=wildfireDb.Collection("updates").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(86400),
	})
	return err
}