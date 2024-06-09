package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"todo/internal/models"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	CreateUserLog(username string, logMessage string) models.Message
}

type service struct {
	db   *mongo.Client
	name string
}

var (
	// host     = os.Getenv("DB_HOST")
	// port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_DATABASE")
)

func New() Service {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	DbName := "todo-app"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(database).SetServerAPIOptions(serverAPI))
	createTodoCollectionIfDoesntExist(client, DbName)

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db:   client,
		name: DbName,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) CreateUserLog(username string, logMessage string) models.Message {
	collection := s.db.Database(s.name).Collection("users_log")
	logBody := models.UserLog{
		Username:   username,
		LogMessage: logMessage,
		CreatedAt:  time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), logBody)
	if err != nil {
		log.Println(err)
		return models.Message{
			Text: "Error Logging user",
		}
	}

	return models.Message{
		Text: "Success Log created",
	}
}

func createTodoCollectionIfDoesntExist(db *mongo.Client, dbName string) {
	databases, err := db.ListDatabaseNames(context.Background(), bson.M{})

	containsDbName := false
	if err != nil {
		log.Fatal(err)
	}

	for _, db := range databases {
		if db == dbName {
			containsDbName = true
		}
	}

	if !containsDbName {
		collection := db.Database(dbName).Collection("users_log")
		_, err := collection.InsertOne(context.Background(), bson.M{
			"system": "createdDb",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("todo-app Db created with todos collection ")
		return
	}
	log.Println("Db todo-app already exists, no need for creation")

}
