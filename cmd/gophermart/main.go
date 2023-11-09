package main

import (
	"crypto/rand"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/interfaces/handlers"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/loggers"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/migrations"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/services"
	"net/http"
	"os"
)

func main() {
	if err := parseFlags(); err != nil {
		fmt.Fprint(os.Stderr, "can't parse flags\n", err)
		os.Exit(1)
	}

	if err := loggers.NewLogger(cfg.LogLevel); err != nil {
		loggers.Logf.Errorf("can't create logger %v", err)
		os.Exit(1)
	}

	err := migrations.Start(cfg.DataBaseURI)
	if err != nil {
		loggers.Logf.Errorf("migration of data to DB is failed %v", err)
		os.Exit(1)
	}

	store, err := repositories.NewStore(cfg.DataBaseURI)
	if err != nil {
		loggers.Logf.Errorf("can't connected to DB %v", err)
		os.Exit(1)
	}

	service := services.NewUserService(store)
	secretKey := make([]byte, 16)
	_, err = rand.Read(secretKey)
	if err != nil {
		loggers.Logf.Errorf("it create secretKey has been finished is wrong %v", err)
		os.Exit(1)
	}
	handlerUser := handlers.NewHandlerUser(service, secretKey)
	router := handlers.UserRouter(handlerUser)
	loggers.Logf.Infof("starting HTTP server on address: %s", cfg.ConnectAddr)

	err = http.ListenAndServe(cfg.ConnectAddr, router)
	if err != nil {
		loggers.Logf.Errorf("error HTTP server %v", err)
		os.Exit(1)
	}
}
