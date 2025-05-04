package handler

import "net/http"

type ClientHandler interface {
	handleClients(w http.ResponseWriter, r *http.Request)
	handleClientByID(w http.ResponseWriter, r *http.Request)
}

type ProxyHandler interface {
	handleProxyRequest(w http.ResponseWriter, r *http.Request)
}
