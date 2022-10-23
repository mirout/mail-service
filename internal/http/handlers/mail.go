package handlers

import (
	"context"
	"encoding/json"
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
	SendMailToGroup(w http.ResponseWriter, r *http.Request)
	GetMailsSentToUser(w http.ResponseWriter, r *http.Request)
	GetMailById(w http.ResponseWriter, r *http.Request)
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
	r.Get("/to/user/{user_id}", s.GetMailsSentToUser)
	r.Get("/{mail_id}", s.GetMailById)
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

	var mail model.MailJson
	err = json.NewDecoder(r.Body).Decode(&mail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mail.Subject == "" || mail.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.sendToUser(r.Context(), mail, user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	var mail model.MailJson
	err = json.NewDecoder(r.Body).Decode(&mail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mail.Subject == "" || mail.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, user := range users {
		err = s.sendToUser(r.Context(), mail, user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *mailHandlers) sendToUser(ctx context.Context, mail model.MailJson, user model.User) error {
	if mail.SendAt == "" {
		return s.sender.SendMailToUser(ctx, model.Mail{
			ToUserId: user.ID,
			Subject:  mail.Subject,
			Body:     mail.Body,
		})
	} else {
		parse, err := time.Parse(time.RFC3339, mail.SendAt)
		if err != nil {
			return err
		}
		return s.sender.CreateDelayedMail(ctx, model.Mail{
			ToUserId: user.ID,
			Subject:  mail.Subject,
			Body:     mail.Body,
		}, parse)
	}
}

func (s *mailHandlers) GetMailsSentToUser(w http.ResponseWriter, r *http.Request) {
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

	mails, err := s.sender.GetMailsBySentTo(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(mails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *mailHandlers) GetMailById(w http.ResponseWriter, r *http.Request) {
	mailId := chi.URLParam(r, "mail_id")
	id, err := uuid.Parse(mailId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mail, err := s.sender.GetMailById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(mail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
