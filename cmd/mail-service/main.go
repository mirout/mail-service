package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"log"
	"mail-service/internal/http"
	"mail-service/internal/storage"
	"os"
)

type Options struct {
	SmtpHost string `long:"smtp-host" description:"SMTP host" required:"true"`
	SmtpPort uint   `long:"smtp-port" description:"SMTP port" required:"true"`

	ServerPort uint `long:"server-port" description:"Server port" default:"8080"`

	DBHost string `long:"db-host" description:"DB host" required:"true"`
	DBPort uint   `long:"db-port" description:"DB port" default:"5432"`
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

	ms := http.NewMailServer(sqlStorage)

	ms.ListenAndServe()

}
