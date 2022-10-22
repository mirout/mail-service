package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"mail-service/internal/model"
	"mail-service/internal/storage"
	"mime"
	"net/http"
)

type MailServer struct {
	*http.Server
	storage storage.StorageService
}

func (s *MailServer) PostCreateUser(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(id.String()))
	log.Print(id)
}

func NewMailServer(storage storage.StorageService) *MailServer {
	s := &MailServer{
		Server: &http.Server{
			Addr: ":8080",
		},
		storage: storage,
	}
	r := chi.NewRouter()
	r.Route("/api/v1/users/create", func(r chi.Router) {
		r.Post("/", s.PostCreateUser)
	})

	s.Handler = r
	return s
}
