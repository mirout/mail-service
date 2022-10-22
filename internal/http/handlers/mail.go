package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mail-service/internal/email"
	"mail-service/internal/storage"
	"net/http"
)

type MailHandlers interface {
	Register(r chi.Router)
	SendMailToUser(w http.ResponseWriter, r *http.Request)
	SendMailToGroup(w http.ResponseWriter, r *http.Request)
}

type mailHandlers struct {
	groups storage.GroupService
	users  storage.UserService
	sender email.MailSender
}

func NewMailHandlers(groups storage.GroupService, users storage.UserService, sender email.MailSender) MailHandlers {
	return &mailHandlers{groups: groups, users: users, sender: sender}
}

func (s *mailHandlers) Register(r chi.Router) {
	r.Post("/to/user/{user_id}", s.SendMailToUser)
	r.Post("/to/group/{group_id}", s.SendMailToGroup)
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

	err = s.sender.SendMailToUser(user.Email, msg.String())
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
		err = s.sender.SendMailToUser(user.Email, msg.String())
	}
}
