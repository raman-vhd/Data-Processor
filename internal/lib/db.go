package lib

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

type Database struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
}

func NewDB(env Env) Database {
	// Mongo
	clientOptions := options.Client().ApplyURI(env.DBURI)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed connecting to database: %v\n", err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("failed pinging database: %v\n", err)
	}

	log.Println("Connected to MongoDB!")

	// Redis
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return Database{
		MongoClient: mongoClient,
        RedisClient: rc,
	}
}

func (d Database) GetCollection(collectionName string) *mongo.Collection {
	collection := d.MongoClient.Database("data_processor").Collection(collectionName)
	return collection
}
