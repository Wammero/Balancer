package handler

import (
	"net/http"

	"github.com/Wammero/Balancer/internal/balancer"
	"github.com/Wammero/Balancer/internal/middleware"
	"github.com/Wammero/Balancer/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	ClientHandler ClientHandler
	ProxyHandler  ProxyHandler
}

func New(service *service.Service, bal *balancer.Balancer, logger *logrus.Logger) *Handler {
	return &Handler{
		ClientHandler: NewClientHandler(service.ClientService, logger),
		ProxyHandler:  NewProxyHandler(service.ProxyService, bal, logger),
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// CRUD clients
	mux.HandleFunc("/clients", h.ClientHandler.handleClients)
	mux.HandleFunc("/clients/", h.ClientHandler.handleClientByID)

	mux.Handle("/", middleware.JWTValidator(http.HandlerFunc(h.ProxyHandler.handleProxyRequest)))
}
