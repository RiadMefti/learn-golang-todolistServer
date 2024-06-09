package models

import "time"

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
	Username   string    `json:"username"`
	Title      string    `json:"title"`
	IsDone     bool      `json:"isDone"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
