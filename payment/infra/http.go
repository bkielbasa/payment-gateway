package infra

import (
	"encoding/json"
	"net/http"
	"payment-gateway/payment/app"
)

type httpHandler struct {
	app app.Payment
}

func NewHandler(app app.Payment) httpHandler {
	return httpHandler{app: app}
}

func responseOK(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func internalError(w http.ResponseWriter, msg string) {
	httpError(w, msg, http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, msg string) {
	httpError(w, msg, http.StatusBadRequest)
}

func notFound(w http.ResponseWriter, msg string) {
	httpError(w, msg, http.StatusNotFound)
}

func unauthorized(w http.ResponseWriter, msg string) {
	httpError(w, msg, http.StatusUnauthorized)
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	type response struct {
		Message string `json:"message"`
	}
	resp := response{
		Message: msg,
	}

	body, _ := json.Marshal(resp)
	w.Write(body)
}

func responseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
