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

type GroupHandlers interface {
	Register(r chi.Router)
	PostCreateGroup(w http.ResponseWriter, r *http.Request)
	GetGroupInfo(w http.ResponseWriter, r *http.Request)
	AddUserToGroup(w http.ResponseWriter, r *http.Request)
	RemoveUserFromGroup(w http.ResponseWriter, r *http.Request)
}

type groupHandlers struct {
	storage storage.GroupService
}

func NewGroupHandlers(storage storage.GroupService) GroupHandlers {
	return &groupHandlers{storage: storage}
}

func (s *groupHandlers) Register(r chi.Router) {
	r.Post("/", s.PostCreateGroup)
	r.Get("/", s.GetGroupInfo)
	r.Post("/{group_id}/add/{user_id}", s.AddUserToGroup)
	r.Post("/{group_id}/remove/{user_id}", s.RemoveUserFromGroup)
}

func (s *groupHandlers) PostCreateGroup(w http.ResponseWriter, r *http.Request) {
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

	var group model.Group
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.storage.CreateGroup(r.Context(), group)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id.String()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *groupHandlers) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	groupId := r.URL.Query().Get("id")
	groupName := r.URL.Query().Get("name")
	if groupId == "" && groupName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var group model.Group
	var err error

	if groupId != "" {
		groupId, err := uuid.Parse(groupId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		group, err = s.storage.GetGroupById(r.Context(), groupId)
	} else {
		group, err = s.storage.GetGroupByName(r.Context(), groupName)
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *groupHandlers) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := uuid.Parse(chi.URLParam(r, "group_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.storage.AddUserToGroup(r.Context(), userId, groupId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *groupHandlers) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := uuid.Parse(chi.URLParam(r, "group_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.storage.RemoveUserFromGroup(r.Context(), userId, groupId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
