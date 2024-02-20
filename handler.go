package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CityBasedService     *Service
	ProvinceBasedService *Service
}

func (h *Handler) GetPredictionPage(c *gin.Context) {
	provinceBasedPrediction, err := h.ProvinceBasedService.GetNationalCountingPrediction(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Maaf ada kesalahan")
		log.Printf("Error getting province-based prediction: %s", err)
		return
	}

	cityBasedPrediction, err := h.CityBasedService.GetNationalCountingPrediction(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Maaf ada kesalahan")
		log.Printf("Error getting city-based prediction: %s", err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"TitleOne":                 "Prediksi Berdasarkan Kota",
		"CandidateOnePercentOne":   fmt.Sprintf("%.2f", cityBasedPrediction.CandidateOnePercent()),
		"CandidateTwoPercentOne":   fmt.Sprintf("%.2f", cityBasedPrediction.CandidateTwoPercent()),
		"CandidateThreePercentOne": fmt.Sprintf("%.2f", cityBasedPrediction.CandidateThreePercent()),
		"UpdatedAtOne":             time.Unix(0, cityBasedPrediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),

		"TitleTwo":                 "Prediksi Berdasarkan Provinsi",
		"CandidateOnePercentTwo":   fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateOnePercent()),
		"CandidateTwoPercentTwo":   fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateTwoPercent()),
		"CandidateThreePercentTwo": fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateThreePercent()),
		"UpdatedAtTwo":             time.Unix(0, provinceBasedPrediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),
	})
}
