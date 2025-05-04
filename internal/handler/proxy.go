package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/Wammero/Balancer/internal/balancer"
	"github.com/Wammero/Balancer/internal/logger"
	"github.com/Wammero/Balancer/internal/service"
	"github.com/Wammero/Balancer/pkg/responsemaker"
	"github.com/sirupsen/logrus"
)

type proxyHandler struct {
	service service.ProxyService
	bal     *balancer.Balancer
	logger  *logrus.Logger
}

func NewProxyHandler(service service.ProxyService, bal *balancer.Balancer, logger *logrus.Logger) *proxyHandler {
	return &proxyHandler{service: service, bal: bal, logger: logger}
}

func (p *proxyHandler) handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reqID := r.Header.Get("X-Request-ID")

	if err := p.service.CheckRateLimit(r.Context()); err != nil {
		logger.LogRequest(p.logger, logrus.WarnLevel, reqID, r.Method, r.URL.Path, "", http.StatusTooManyRequests, time.Since(start), err)
		responsemaker.WriteJSONError(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	backend := p.bal.GetNextBackend()
	if backend == nil {
		logger.LogRequest(p.logger, logrus.WarnLevel, reqID, r.Method, r.URL.Path, "", http.StatusServiceUnavailable, time.Since(start), nil)
		responsemaker.WriteJSONError(w, "No alive backends available", http.StatusServiceUnavailable)
		return
	}

	req, _ := http.NewRequest(r.Method, backend.URL.String()+r.URL.Path, r.Body)
	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		backend.SetAlive(false)
		logger.LogRequest(p.logger, logrus.ErrorLevel, reqID, r.Method, r.URL.Path, backend.URL.String(), http.StatusBadGateway, time.Since(start), err)
		responsemaker.WriteJSONError(w, "Failed to reach backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	logger.LogRequest(p.logger, logrus.InfoLevel, reqID, r.Method, r.URL.Path, backend.URL.String(), resp.StatusCode, time.Since(start), nil)
}
