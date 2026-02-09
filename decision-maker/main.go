package main

import (
	"decisionMaker/api"
	"decisionMaker/persistence/postgres"
	"decisionMaker/service/device"
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
	fmt.Println(Banner())

	db, err := postgres.GetInstance()
	if err != nil {
		log.Fatal("error: connecting to database", err)
	}

	// dep injections
	deviceRepository := postgres.NewDeviceRepository(db)
	deviceService := device.NewService(deviceRepository)

	server := api.NewServer(deviceService)

	router := chi.NewMux()
	router.Use(middleware.Logger)
	si := api.NewStrictHandler(server, nil)
	h := api.HandlerFromMuxWithBaseURL(si, router, "/api")

	log.Println("ready to accept requests")

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:3000",
	}
	log.Fatal(s.ListenAndServe())
}
