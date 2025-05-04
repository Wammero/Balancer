package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Wammero/Balancer/internal/logger"
	"github.com/Wammero/Balancer/internal/models"
	"github.com/Wammero/Balancer/internal/service"
	"github.com/Wammero/Balancer/pkg/responsemaker"
	"github.com/sirupsen/logrus"
)

type clientHandler struct {
	service service.ClientService
	logger  *logrus.Logger
}

func NewClientHandler(service service.ClientService, logger *logrus.Logger) *clientHandler {
	return &clientHandler{service: service, logger: logger}
}

func (c *clientHandler) handleClients(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reqID := r.Header.Get("X-Request-ID")

	logger.LogRequest(c.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, "", http.StatusOK, time.Since(start), nil)

	switch r.Method {
	case http.MethodPost:
		// Добавление нового клиента

		var req models.Client
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, "", http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		token, err := c.service.CreateClient(r.Context(), req.ClientID, req.Capacity, req.RatePerSec)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, "", http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Failed to create client", http.StatusBadRequest)
			return
		}

		logger.LogRequest(logrus.StandardLogger(), logrus.InfoLevel, reqID, r.Method, r.URL.Path, "", http.StatusCreated, time.Since(start), nil)
		responsemaker.WriteJSONResponse(w, token, http.StatusCreated)

	case http.MethodGet:
		// Получение списка всех клиентов

		Clients, err := c.service.ListClients(r.Context())
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, "", http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Failed to get clients", http.StatusBadRequest)
			return
		}

		logger.LogRequest(c.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, "", http.StatusOK, time.Since(start), nil)
		responsemaker.WriteJSONResponse(w, Clients, http.StatusOK)

	default:
		responsemaker.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		logger.LogRequest(c.logger, logrus.WarnLevel, reqID, r.Method, r.URL.Path, "", http.StatusMethodNotAllowed, time.Since(start), nil)
	}
}

func (c *clientHandler) handleClientByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reqID := r.Header.Get("X-Request-ID")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, "", http.StatusBadRequest, time.Since(start), fmt.Errorf("Invalid URL format"))
		responsemaker.WriteJSONError(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	clientID := parts[2]

	switch r.Method {
	case http.MethodGet:
		// Получить одного клиента по ID

		client, err := c.service.GetClientByID(r.Context(), clientID)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Failed to get client", http.StatusBadRequest)
			return
		}

		logger.LogRequest(c.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusOK, time.Since(start), nil)
		responsemaker.WriteJSONResponse(w, client, http.StatusOK)

	case http.MethodDelete:
		// Удалить клиента по ID

		err := c.service.DeleteClient(r.Context(), clientID)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Failed to delete client", http.StatusBadRequest)
			return
		}

		logger.LogRequest(c.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusOK, time.Since(start), nil)
		responsemaker.WriteJSONResponse(w, "OK", http.StatusOK)

	case http.MethodPut:
		// Обновить лимиты клиента

		var req models.Client
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusBadRequest, time.Since(start), err)
			responsemaker.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = c.service.UpdateClient(r.Context(), clientID, req.Capacity, req.RatePerSec)
		if err != nil {
			logger.LogRequest(c.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusBadRequest, time.Since(start), err)

			return
		}

		logger.LogRequest(c.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusOK, time.Since(start), nil)
		responsemaker.WriteJSONResponse(w, "OK", http.StatusOK)

	default:
		responsemaker.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		logger.LogRequest(c.logger, logrus.WarnLevel, reqID, r.Method, r.URL.Path, clientID, http.StatusMethodNotAllowed, time.Since(start), nil)
	}
}
