package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"log"
	"mail-service/internal/queue"
	"mail-service/internal/services"
	"mail-service/internal/services/group"
	"mail-service/internal/services/img"
	"mail-service/internal/services/mail"
	"mail-service/internal/services/user"
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
	MailHost     string `long:"mail-host" description:"Mail host" required:"true"`

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

	delayedQueue, err := queue.NewQueue(context.Background(), redisAddr, os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		log.Fatalf("Can't create delayedQueue: %v", err)
	}
	go delayedQueue.Run()

	mailServerAddr := fmt.Sprintf("%s:%d", opts.SmtpHost, opts.SmtpPort)
	smtpConf := mail.SmtpConfig{Addr: mailServerAddr, Username: opts.MailUsername, Password: opts.MailPassword}

	mailSender, err := mail.NewWorker(opts.MailHost, smtpConf, sqlStorage, sqlStorage, delayedQueue)
	if err != nil {
		log.Fatalf("Can't create mail server: %v", err)
	}
	defer func(mailSender *mail.MailWorker) {
		err := mailSender.Close()
		if err != nil {
			log.Printf("Can't close mail server: %v", err)
		}
	}(mailSender)

	go mailSender.Run()

	h := services.NewMailServer(
		user.NewUserHandlers(sqlStorage),
		group.NewGroupHandlers(sqlStorage),
		mail.NewMailHandlers(sqlStorage, sqlStorage, mailSender),
		img.NewImageHandlers(sqlStorage),
		opts.ServerPort,
	)

	if err = h.ListenAndServe(); err != nil {
		log.Fatalf("Can't start server: %v", err)
	}

}
