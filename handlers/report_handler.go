package handlers

import (
	"category-api/services"
	"encoding/json"
	"net/http"
)

type ReportHandler struct {
	sevice *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{sevice: service}
}

func (h *ReportHandler) getSellingToday(w http.ResponseWriter, r *http.Request) {
	sellingToday, err := h.sevice.GetSalesToday()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sellingToday)
}

func (h *ReportHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSellingToday(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
