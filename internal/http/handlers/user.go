package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"mail-service/internal/model"
	"mail-service/internal/storage"
	"mime"
	"net/http"
)

type UserHandlers interface {
	Register(r chi.Router)
	PostCreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
}

type userHandlers struct {
	storage storage.UserService
}

func NewUserHandlers(storage storage.UserService) UserHandlers {
	return &userHandlers{storage: storage}
}

func (s *userHandlers) Register(r chi.Router) {
	r.Post("/", s.PostCreateUser)
	r.Get("/", s.GetUser)
}

func (s *userHandlers) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	t, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if t != "application/json" {
		log.Println(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.storage.CreateUser(r.Context(), user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id.String()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *userHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	email := r.URL.Query().Get("email")
	if userId == "" && email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	var err error

	if userId != "" {
		id, err := uuid.Parse(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err = s.storage.GetUser(r.Context(), id)
	} else if email != "" {
		user, err = s.storage.GetUserByEmail(r.Context(), email)
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
