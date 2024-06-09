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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	CreateUserLog(username string, logMessage string) models.Message
	GetUserTodos(username string) ([]models.Todo, error)
	PostTodo(username string, title string) ([]models.Todo, error)
	PutTodo(id string, username string, isDone bool) ([]models.Todo, error) 
	DeleteTodo(id string, username string) ([]models.Todo, error)
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

func (s *service) GetUserTodos(username string) ([]models.Todo, error) {

	var todos []models.Todo
	filter := bson.M{"username": username}

	collection := s.db.Database(s.name).Collection("todos")

	dbTodos, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = dbTodos.All(context.Background(), &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil

}

func (s *service) PostTodo(username string, title string) ([]models.Todo, error) {
	var todo = models.Todo{
		Username:   username,
		Title:      title,
		IsDone:     false,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	collection := s.db.Database(s.name).Collection("todos")
	_, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	return s.GetUserTodos(username)

}
func (s *service) PutTodo(id string, username string, isDone bool) ([]models.Todo, error) {
	// Convert the id string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"isDone": isDone}}
	collection := s.db.Database(s.name).Collection("todos")

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetUserTodos(username)
}
func (s *service) DeleteTodo(id string, username string) ([]models.Todo, error) {
	// Convert the id string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	collection := s.db.Database(s.name).Collection("todos")

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return s.GetUserTodos(username)
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
