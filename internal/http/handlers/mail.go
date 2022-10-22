package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"mail-service/internal/email"
	"mail-service/internal/model"
	"mail-service/internal/storage"
	"net/http"
	"time"
)

type MailHandlers interface {
	Register(r chi.Router)
	SendMailToUser(w http.ResponseWriter, r *http.Request)
	SendDelayedMailToUser(w http.ResponseWriter, r *http.Request)
	SendMailToGroup(w http.ResponseWriter, r *http.Request)
	SendDelayedMailToGroup(w http.ResponseWriter, r *http.Request)
}

type mailHandlers struct {
	groups storage.Group
	users  storage.User
	sender email.MailSender
}

func NewMailHandlers(groups storage.Group, users storage.User, sender email.MailSender) MailHandlers {
	return &mailHandlers{groups: groups, users: users, sender: sender}
}

func (s *mailHandlers) Register(r chi.Router) {
	r.Post("/to/user/{user_id}", s.SendMailToUser)
	r.Post("/to/group/{group_id}", s.SendMailToGroup)
	r.Post("/to/user/{user_id}/{at}", s.SendDelayedMailToUser)
	r.Post("/to/group/{group_id}/{at}", s.SendDelayedMailToGroup)
}

func (s *mailHandlers) SendMailToUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.GetUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := bytes.Buffer{}
	_, err = msg.ReadFrom(r.Body)

	err = s.sender.SendMailToUser(r.Context(), model.Mail{
		ToUserId: user.ID,
		Subject:  "Test",
		Body:     msg.String(),
	})
}

func (s *mailHandlers) SendMailToGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
	id, err := uuid.Parse(groupId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := s.groups.GetUsersByGroup(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := bytes.Buffer{}
	_, err = msg.ReadFrom(r.Body)

	for _, user := range users {
		err = s.sender.SendMailToUser(r.Context(), model.Mail{
			ToUserId: user.ID,
			Subject:  "Test",
			Body:     msg.String(),
		})
	}
}

func (s *mailHandlers) SendDelayedMailToUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
	id, err := uuid.Parse(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	at := chi.URLParam(r, "at")
	parse, err := time.Parse(time.RFC3339, at)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.users.GetUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := bytes.Buffer{}
	_, err = msg.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.sender.CreateDelayedMail(r.Context(), model.Mail{
		ToUserId: user.ID,
		Subject:  "Test",
		Body:     msg.String(),
	}, parse)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *mailHandlers) SendDelayedMailToGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
	id, err := uuid.Parse(groupId)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	at := chi.URLParam(r, "at")
	parse, err := time.Parse(time.RFC3339, at)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := s.groups.GetUsersByGroup(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := bytes.Buffer{}
	_, err = msg.ReadFrom(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		err = s.sender.CreateDelayedMail(r.Context(), model.Mail{
			ToUserId: user.ID,
			Subject:  "Test",
			Body:     msg.String(),
		}, parse)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
