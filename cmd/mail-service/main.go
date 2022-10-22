package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"log"
	"mail-service/internal/email"
	"mail-service/internal/http"
	"mail-service/internal/http/handlers"
	"mail-service/internal/storage"
	"os"
)

type Options struct {
	SmtpHost string `long:"smtp-host" description:"SMTP host" required:"true"`
	SmtpPort uint   `long:"smtp-port" description:"SMTP port" required:"true"`

	ServerPort int `long:"server-port" description:"Server port" default:"8080"`

	DBHost string `long:"db-host" description:"DB host" required:"true"`
	DBPort uint   `long:"db-port" description:"DB port" default:"5432"`

	MailUsername string `long:"mail-username" description:"Mail username" required:"true"`
	MailPassword string `long:"mail-password" description:"Mail password" required:"true"`
}

var appName = "mail-service"

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Can't parse options")
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.DBHost,
		opts.DBPort,
		appName,
		os.Getenv("DB_PASSWORD"),
		appName,
	)

	sqlStorage, err := storage.NewSqlStorage(context.Background(), "postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Can't create sql storage: %v", err)
	}

	mailServerAddr := fmt.Sprintf("%s:%d", opts.SmtpHost, opts.SmtpPort)
	mailServer, err := email.NewSender(mailServerAddr, opts.MailUsername, opts.MailPassword)

	if err != nil {
		log.Fatalf("Can't create mail server: %v", err)
	}

	handlers := http.NewMailServer(
		handlers.NewUserHandlers(sqlStorage),
		handlers.NewGroupHandlers(sqlStorage),
		handlers.NewMailHandlers(sqlStorage, sqlStorage, mailServer),
		opts.ServerPort,
	)

	handlers.ListenAndServe()

}
