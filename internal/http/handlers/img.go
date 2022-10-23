package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"mail-service/internal/storage"
	"net/http"
	"os"
)

type ImageHandlers interface {
	Register(r chi.Router)
	GetImage(w http.ResponseWriter, r *http.Request)
}

type imageHandlers struct {
	mails storage.Mail
}

func NewImageHandlers(mails storage.Mail) ImageHandlers {
	return &imageHandlers{mails: mails}
}

func (s *imageHandlers) Register(r chi.Router) {
	r.Get("/{mail_id}.png", s.GetImage)
}

func (s *imageHandlers) GetImage(w http.ResponseWriter, r *http.Request) {
	mailId := chi.URLParam(r, "mail_id")
	id, err := uuid.Parse(mailId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.mails.MarkAsWatched(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := os.Open("templates/img.png")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, file)
}
