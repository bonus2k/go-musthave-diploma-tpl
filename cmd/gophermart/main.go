package main

import (
	"crypto/rand"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/interfaces/clients"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/interfaces/handlers"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/migrations"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/services"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := parseFlags(); err != nil {
		fmt.Fprint(os.Stderr, "can't parse flags\n", err)
		os.Exit(1)
	}

	if err := internal.NewLogger(cfg.LogLevel); err != nil {
		internal.Logf.Errorf("can't create logger %v", err)
		os.Exit(1)
	}

	err := migrations.Start(cfg.DataBaseURI)
	if err != nil {
		internal.Logf.Errorf("migration of data to DB is failed %v", err)
		os.Exit(1)
	}

	store, err := repositories.NewStore(cfg.DataBaseURI)
	if err != nil {
		internal.Logf.Errorf("can't connected to DB %v", err)
		os.Exit(1)
	}

	service := services.NewUserService(store)
	secretKey := make([]byte, 16)
	_, err = rand.Read(secretKey)
	if err != nil {
		internal.Logf.Errorf("it create secretKey has been finished is wrong %v", err)
		os.Exit(1)
	}

	client := resty.New()
	accrual := clients.NewClientAccrual(client, "http://localhost:8081")
	worker := clients.NewPoolWorker(accrual, service)
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		worker.StarIntegration(5)
		time.Sleep(1 * time.Minute)
		close(worker.Done)
	}()

	go func() {
		for range ticker.C {
			worker.StarIntegration(5)
			time.Sleep(1 * time.Minute)
			close(worker.Done)
		}
	}()

	//resp, err := client.R().
	//	SetPathParams(map[string]string{"g": "get", "url": "httpbin"}).
	//	EnableTrace().
	//	Get("https://{url}.org/{g}")
	//if err != nil {
	//	fmt.Println(err)
	//}

	handlerUser := handlers.NewHandlerUser(service, secretKey)
	router := handlers.UserRouter(handlerUser)
	internal.Logf.Infof("starting HTTP server on address: %s", cfg.ConnectAddr)
	err = http.ListenAndServe(cfg.ConnectAddr, router)
	if err != nil {
		internal.Logf.Errorf("error HTTP server %v", err)
		os.Exit(1)
	}

}
