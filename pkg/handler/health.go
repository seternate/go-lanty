package handler

import (
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (handler *HealthHandler) GetHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
