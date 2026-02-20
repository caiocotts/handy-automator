package main

import (
	"decisionMaker/api"
	"decisionMaker/persistence/postgres"
	"decisionMaker/service/device"
	"decisionMaker/service/workflow"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("error: loading .env:", err)
	}
}

func main() {
	const baseUrl = "/api"

	fmt.Println(Banner())

	db, err := postgres.GetInstance()
	if err != nil {
		log.Fatal("error: connecting to database", err)
	}

	// dep injections
	deviceRepository := postgres.NewDeviceRepository(db)
	workflowRepository := postgres.NewWorkflowRepository(db)

	deviceService := device.NewService(deviceRepository)
	workflowService := workflow.NewService(workflowRepository)

	server := api.NewServer(deviceService, workflowService)

	router := chi.NewMux()
	router.Use(middleware.Logger)

	router.Get(baseUrl+"/docs", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(api.Docs)
		if err != nil {
			log.Print(err)
		}
	})

	router.Get(baseUrl+"/spec", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(api.Spec)
		if err != nil {
			log.Print(err)
		}
	})

	si := api.NewStrictHandler(server, nil)
	h := api.HandlerFromMuxWithBaseURL(si, router, baseUrl)

	log.Println("ready to accept requests")

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:3000",
	}
	log.Fatal(s.ListenAndServe())
}
