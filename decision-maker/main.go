package main

import (
	"context"
	"decisionMaker/api"
	"decisionMaker/config"
	"decisionMaker/consts"
	"decisionMaker/persistence/postgres"
	"decisionMaker/service/auth"
	"decisionMaker/service/device"
	"decisionMaker/service/discovery"
	"decisionMaker/service/user"
	"decisionMaker/service/workflow"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.Load()
	fmt.Println(Banner())

	db, err := postgres.GetInstance()
	if err != nil {
		log.Fatal("error: unable to connect to database: ", err)
	}

	// dep injections
	deviceRepository := postgres.NewDeviceRepository(db)
	userRepository := postgres.NewUserRepository(db)
	workflowRepository := postgres.NewWorkflowRepository(db)

	discoveryService := discovery.NewService(deviceRepository)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go discoveryService.Start(ctx)

	deviceService := device.NewService(deviceRepository)
	userService := user.NewService(userRepository)
	workflowService := workflow.NewService(workflowRepository, discoveryService)
	authService := auth.NewService(userRepository)

	server := api.NewServer(deviceService, userService, workflowService, authService, discoveryService)

	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(api.JWTAuthMiddleware(authService.ValidateToken))

	router.Get(consts.APIBaseUrl+"/docs", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(api.Docs)
		if err != nil {
			log.Print(err)
		}
	})

	router.Get(consts.APIBaseUrl+"/spec", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(api.Spec)
		if err != nil {
			log.Print(err)
		}
	})

	si := api.NewStrictHandler(server, nil)
	h := api.HandlerFromMuxWithBaseURL(si, router, consts.APIBaseUrl)

	log.Println("ready to accept requests")

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:3000",
	}
	log.Fatal(s.ListenAndServe())
}
