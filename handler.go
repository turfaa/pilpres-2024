package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetPredictionPage(c *gin.Context) {
	prediction, err := h.Service.GetNationalCountingPrediction(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Maaf ada kesalahan")
		log.Printf("Error getting prediction: %s", err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"CandidateOnePercent":   fmt.Sprintf("%.2f", prediction.CandidateOnePercent()),
		"CandidateTwoPercent":   fmt.Sprintf("%.2f", prediction.CandidateTwoPercent()),
		"CandidateThreePercent": fmt.Sprintf("%.2f", prediction.CandidateThreePercent()),
		"UpdatedAt":             time.Unix(0, prediction.UpdatedAt*int64(time.Millisecond)).Format("02-Jan-2006 15:04:05"),
	})
}
