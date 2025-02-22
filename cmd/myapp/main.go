package main

import (
	"cofee-shop-mongo/pkg/config"
	"cofee-shop-mongo/pkg/lib/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	ConnectionString := cfg.MakeConnectionString()
	fmt.Println("Connection String: ", ConnectionString)
	logger := logger.SetupPrettySlog(os.Stdout)

	client := mongoConnect(ConnectionString)
	defer mongoDisconnect(client)

	mux := http.NewServeMux()
	server := NewAPIServer("localhost:8000", mux, client.Database("cofee-shop"), logger)
	server.Run()
}

func mongoConnect(connectionString string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().
		ApplyURI(connectionString))
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB")
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't ping MongoDB")
	}
	return client
}

func mongoDisconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal("Couldn't disconnect from MongoDB")
	}
}
