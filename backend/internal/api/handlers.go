package api

import (
	"encoding/json"
	"github.com/Rastler3D/container-monitoring/backend/internal/service"
	"github.com/Rastler3D/container-monitoring/common/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Handler struct {
	pingService *service.PingService
	logger      *log.Logger
}

func NewHandler(pingService *service.PingService, logger *log.Logger) Handler {
	return Handler{pingService: pingService, logger: logger}
}

func (h *Handler) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/containers", h.getContainers).Methods("GET")
	r.HandleFunc("/api/containers", h.addContainers).Methods("POST")
	return r
}

func (h *Handler) getContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := h.pingService.GetAllContainerStatuses()
	if err != nil {
		h.logger.Printf("Error getting containers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}

func (h *Handler) addContainers(w http.ResponseWriter, r *http.Request) {
	var c []model.ContainerStatus
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.pingService.AddContainerStatuses(c); err != nil {
		h.logger.Printf("Error adding container: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
