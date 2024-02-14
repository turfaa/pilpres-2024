package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	service := &Service{
		KawalPemiluClient: &KawalPemiluClient{
			BaseURL: "https://kp24-fd486.et.r.appspot.com",
		},
		Predictor: &SimplePredictor{},
	}

	go service.RunRefresher(context.Background())

	handler := &Handler{Service: service}

	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	router.StaticFile("/favicon.ico", "static/favicon.ico")

	router.GET("/", handler.GetPredictionPage)

	errChan := make(chan error)
	go func() {
		errChan <- router.Run(":8080")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGINT)

	select {
	case err := <-errChan:
		if err != nil {
			log.Fatal(err)
		}

	case <-done:
		log.Println("shutdown")
	}
}
