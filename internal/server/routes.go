package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	r.Post("/log", s.createLog)

	r.Post("/todo", s.createTodo)
	r.Get("/todos", s.getTodos)

	r.Put("/todo", s.modifyTodo)
	r.Delete("/todo", s.deleteTodo)
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) createLog(w http.ResponseWriter, r *http.Request) {
	var logReq models.LogRequest

	err := json.NewDecoder(r.Body).Decode(&logReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResp, _ := json.Marshal(s.db.CreateUserLog(logReq.Username, logReq.LogMessage))
	_, _ = w.Write(jsonResp)

}
func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var todoReq models.TodoRequest

	err := json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todos, err := s.db.PostTodo(todoReq.Username, todoReq.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResp, _ := json.Marshal(todos)
	_, _ = w.Write(jsonResp)
}
func (s *Server) getTodos(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	username := queryParams.Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	todos, err := s.db.GetUserTodos(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResp, _ := json.Marshal(todos)
	_, _ = w.Write(jsonResp)
}

func (s *Server) modifyTodo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")
	username := queryParams.Get("username")
	isDoneStr := queryParams.Get("isDone")
	if username == "" || id == "" || isDoneStr == "" {
		http.Error(w, "uall fields are required", http.StatusBadRequest)
		return
	}
	isDone, err := strconv.ParseBool(isDoneStr)
	if err != nil {
		http.Error(w, "isDone must be a boolean", http.StatusBadRequest)
		return
	}

	todos, err := s.db.PutTodo(id, username, isDone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResp, _ := json.Marshal(todos)
	_, _ = w.Write(jsonResp)
}
func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")
	username := queryParams.Get("username")

	if id == "" || username == "" {
		http.Error(w, "uall fields are required", http.StatusBadRequest)
		return
	}

	todos, err := s.db.DeleteTodo(id, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResp, _ := json.Marshal(todos)
	_, _ = w.Write(jsonResp)
}
