package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"maps.patio.com/configuration"
	"maps.patio.com/repository"
	routes "maps.patio.com/routes"
)

func main() {
	ctx := context.Background()
	serverDoneChan := make(chan os.Signal, 1)
	signal.Notify(serverDoneChan, os.Interrupt, syscall.SIGTERM)

	config, err := configuration.New()
	if err != nil {
		log.Fatal(err)
	}

	mMap, err := repository.New(config)
	if err != nil {
		log.Fatal(err)
	}

	port := fmt.Sprintf(":%d", config.APP.Port)
	router := routes.Maps(mMap)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Server started on port " + port)
	<-serverDoneChan
	srv.Shutdown(ctx)
	log.Println("Server stopped")
}
