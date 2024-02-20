package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	kpClient := &KawalPemiluClient{
		BaseURL:     "https://kp24-fd486.et.r.appspot.com",
		Concurrency: 20,
	}

	kecamatanBasedService := &Service{
		CountingResultGetter: kpClient.GetAllCitiesCountingResult,
		Predictor:            &SimplePredictor{},
		RefreshInterval:      20 * time.Minute,
	}

	go kecamatanBasedService.RunRefresher(context.Background())

	cityBasedService := &Service{
		CountingResultGetter: kpClient.GetAllProvincesCountingResult,
		Predictor:            &SimplePredictor{},
		RefreshInterval:      5 * time.Minute,
	}

	go cityBasedService.RunRefresher(context.Background())

	provinceBasedService := &Service{
		CountingResultGetter: kpClient.GetNationalCountingResult,
		Predictor:            &SimplePredictor{},
		RefreshInterval:      time.Minute,
	}

	go provinceBasedService.RunRefresher(context.Background())

	handler := &Handler{
		KecamatanBasedService: kecamatanBasedService,
		CityBasedService:      cityBasedService,
		ProvinceBasedService:  provinceBasedService,
	}

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
