package services

import (
	"github.com/go-chi/chi/v5"
	"log"
	"mail-service/internal/services/group"
	"mail-service/internal/services/img"
	"mail-service/internal/services/mail"
	"mail-service/internal/services/user"
	"net/http"
	"strconv"
)

type MailServer struct {
	*http.Server
	users  user.UserHandlers
	groups group.GroupHandlers
	mails  mail.MailHandlers
	imgs   img.ImageHandlers
}

func NewMailServer(userServer user.UserHandlers, groupServer group.GroupHandlers, mails mail.MailHandlers, imgs img.ImageHandlers, port int) *MailServer {
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
