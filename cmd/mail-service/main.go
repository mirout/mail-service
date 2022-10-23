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

	RedisHost string `long:"redis-host" description:"Redis address" required:"true"`
	RedisPort uint   `long:"redis-port" description:"Redis port" default:"6379"`
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
		os.Getenv("POSTGRES_PASSWORD"),
		appName,
	)

	sqlStorage, err := storage.NewSqlStorage(context.Background(), "postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Can't create sql storage: %v", err)
	}
	defer func(sqlStorage *storage.SqlStorage) {
		err := sqlStorage.Close()
		if err != nil {
			log.Fatalf("Can't close sql storage: %v", err)
		}
	}(sqlStorage)

	redisAddr := fmt.Sprintf("%s:%d", opts.RedisHost, opts.RedisPort)

	queue, err := email.NewQueue(context.Background(), redisAddr, os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		log.Fatalf("Can't create queue: %v", err)
	}
	go queue.Run()

	mailServerAddr := fmt.Sprintf("%s:%d", opts.SmtpHost, opts.SmtpPort)
	smtpConf := email.SmtpConfig{Addr: mailServerAddr, Username: opts.MailUsername, Password: opts.MailPassword}

	mailSender, err := email.NewWorker(smtpConf, sqlStorage, sqlStorage, queue)
	if err != nil {
		log.Fatalf("Can't create mail server: %v", err)
	}
	defer func(mailSender *email.MailWorker) {
		err := mailSender.Close()
		if err != nil {
			log.Printf("Can't close mail server: %v", err)
		}
	}(mailSender)

	go mailSender.Run()

	h := http.NewMailServer(
		handlers.NewUserHandlers(sqlStorage),
		handlers.NewGroupHandlers(sqlStorage),
		handlers.NewMailHandlers(sqlStorage, sqlStorage, mailSender),
		handlers.NewImageHandlers(sqlStorage),
		opts.ServerPort,
	)

	if err = h.ListenAndServe(); err != nil {
		log.Fatalf("Can't start server: %v", err)
	}

}
