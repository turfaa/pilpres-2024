package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	KecamatanBasedService *Service
	CityBasedService      *Service
	ProvinceBasedService  *Service
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

	kecamatanBasedPrediction, err := h.KecamatanBasedService.GetNationalCountingPrediction(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Maaf ada kesalahan")
		log.Printf("Error getting kecamatan-based prediction: %s", err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"TitleOne":                 "Prediksi Berdasarkan Kecamatan",
		"CandidateOnePercentOne":   fmt.Sprintf("%.2f", kecamatanBasedPrediction.CandidateOnePercent()),
		"CandidateTwoPercentOne":   fmt.Sprintf("%.2f", kecamatanBasedPrediction.CandidateTwoPercent()),
		"CandidateThreePercentOne": fmt.Sprintf("%.2f", kecamatanBasedPrediction.CandidateThreePercent()),
		"UpdatedAtOne":             time.Unix(0, kecamatanBasedPrediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),

		"TitleTwo":                 "Prediksi Berdasarkan Kota",
		"CandidateOnePercentTwo":   fmt.Sprintf("%.2f", cityBasedPrediction.CandidateOnePercent()),
		"CandidateTwoPercentTwo":   fmt.Sprintf("%.2f", cityBasedPrediction.CandidateTwoPercent()),
		"CandidateThreePercentTwo": fmt.Sprintf("%.2f", cityBasedPrediction.CandidateThreePercent()),
		"UpdatedAtTwo":             time.Unix(0, cityBasedPrediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),

		"TitleThree":                 "Prediksi Berdasarkan Provinsi",
		"CandidateOnePercentThree":   fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateOnePercent()),
		"CandidateTwoPercentThree":   fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateTwoPercent()),
		"CandidateThreePercentThree": fmt.Sprintf("%.2f", provinceBasedPrediction.CandidateThreePercent()),
		"UpdatedAtThree":             time.Unix(0, provinceBasedPrediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),
	})
}
