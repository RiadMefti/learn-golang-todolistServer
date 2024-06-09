package models

import "time"

type Message struct {
	Text string `json:"message"`
}

type LogRequest struct {
	Username   string `json:"username"`
	LogMessage string `json:"logMessage"`
}

type UserLog struct {
    Username   string    `json:"username"`
    LogMessage string    `json:"logMessage"`
    CreatedAt  time.Time `json:"createdAt"`
}
