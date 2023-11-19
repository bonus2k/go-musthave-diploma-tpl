package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/interfaces/clients"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/interfaces/handlers"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/migrations"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/services"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"time"
)

const (
	countWorker             = 5
	retryTimeCheckNewOrders = 5 * time.Second
	migrationsPath          = "file://internal/migrations/sql"
)

func main() {
	fmt.Fprintln(os.Stdout, "starting server...")
	if err := parseFlags(); err != nil {
		fmt.Fprint(os.Stderr, "can't parse flags\n", err)
		os.Exit(1)
	}

	if err := internal.InitLogger(cfg.LogLevel); err != nil {
		internal.Logf.Errorf("can't create logger %v", err)
		os.Exit(1)
	}

	err := migrations.Start(cfg.DataBaseURI, migrationsPath)
	if err != nil {
		internal.Logf.Errorf("migration of data to DB is failed %v", err)
		os.Exit(1)
	}

	db, err := sqlx.Open("pgx", cfg.DataBaseURI)
	if err != nil {
		internal.Logf.Errorf("can't connected to DB %v", err)
		os.Exit(1)
	}
	store := repositories.NewStore(db)

	service := services.NewUserService(store)
	secretKey, err := base64.StdEncoding.DecodeString(cfg.SecretKey)
	if err != nil || len(secretKey) < 16 {
		secretKey = make([]byte, 16)
		_, err = rand.Read(secretKey)
		if err != nil {
			internal.Logf.Errorf("it create secretKey has been finished is wrong %v", err)
			os.Exit(1)
		}
	}

	internal.Logf.Infof("starting integration to: %s", cfg.AccrualURI)
	client := resty.New()
	accrual := clients.NewClientAccrual(client, cfg.AccrualURI)
	ticker := time.NewTicker(retryTimeCheckNewOrders)
	worker := services.NewPoolWorker(accrual, service)
	go func() {
		worker.StarIntegration(countWorker, ticker)
	}()

	internal.Logf.Infof("starting HTTP server on address: %s", cfg.ConnectAddr)
	handlerUser := handlers.NewHandlerUser(service, secretKey)
	router := handlers.UserRouter(handlerUser)
	err = http.ListenAndServe(cfg.ConnectAddr, router)
	if err != nil {
		internal.Logf.Errorf("error HTTP server %v", err)
		os.Exit(1)
	}

}
