package http

import (
	"github.com/go-chi/chi/v5"
	"mail-service/internal/http/handlers"
	"net/http"
	"strconv"
)

type MailServer struct {
	*http.Server
	users  handlers.UserHandlers
	groups handlers.GroupHandlers
	mails  handlers.MailHandlers
}

func NewMailServer(userServer handlers.UserHandlers, groupServer handlers.GroupHandlers, mails handlers.MailHandlers, port int) *MailServer {
	s := &MailServer{
		Server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
		users:  userServer,
		groups: groupServer,
		mails:  mails,
	}

	r := chi.NewRouter()

	r.Route("/api/v1/users", s.users.Register)
	r.Route("/api/v1/groups", s.groups.Register)
	r.Route("/api/v1/mails", s.mails.Register)

	s.Handler = r
	return s
}
