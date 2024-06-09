package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Text string `json:"message"`
}

type LogRequest struct {
	Username   string `json:"username"`
	LogMessage string `json:"logMessage"`
}


type TodoRequest struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	IsDone   bool   `json:"isDone"`
}

type UserLog struct {
	Username   string    `json:"username"`
	LogMessage string    `json:"logMessage"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Todo struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username   string             `json:"username"`
	Title      string             `json:"title"`
	IsDone     bool               `json:"isDone"`
	CreatedAt  time.Time          `json:"createdAt"`
	ModifiedAt time.Time          `json:"modifiedAt"`
}
