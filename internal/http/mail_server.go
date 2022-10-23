package http

import (
	"github.com/go-chi/chi/v5"
	"log"
	"mail-service/internal/http/handlers"
	"net/http"
	"strconv"
)

type MailServer struct {
	*http.Server
	users  handlers.UserHandlers
	groups handlers.GroupHandlers
	mails  handlers.MailHandlers
	imgs   handlers.ImageHandlers
}

func NewMailServer(userServer handlers.UserHandlers, groupServer handlers.GroupHandlers, mails handlers.MailHandlers, imgs handlers.ImageHandlers, port int) *MailServer {
	s := &MailServer{
		Server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
		users:  userServer,
		groups: groupServer,
		mails:  mails,
		imgs:   imgs,
	}

	r := chi.NewRouter()

	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL)
			handler.ServeHTTP(w, r)
		})
	})

	r.Route("/api/v1/users", s.users.Register)
	r.Route("/api/v1/groups", s.groups.Register)
	r.Route("/api/v1/mails", s.mails.Register)
	r.Route("/img", s.imgs.Register)

	s.Handler = r
	return s
}
