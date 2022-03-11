package infra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"payment-gateway/payment/domain"
)

type captureRequest struct {
	UID    string `json:"uid"`
	Amount int    `json:"amount"`
}

func (h *httpHandler) Capture(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, "failed to read body")
		return
	}

	req := captureRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		badRequest(w, "failed to parse request")
		return
	}

	ctx := r.Context()

	m, err := h.app.Capture(ctx, req.UID, req.Amount)
	if errors.Is(err, domain.ErrTransactionNotEnoughMoney) {
		badRequest(w, "failed to capture: "+err.Error())
		return
	}

	if errors.Is(err, domain.ErrTransactionNotFound) {
		notFound(w, "failed to capture: "+err.Error())
		return
	}

	if errors.Is(err, domain.ErrTransactionUnathorized) {
		unauthorized(w, "failed to capture: "+err.Error())
		return
	}

	if err != nil {
		internalError(w, "failed to authorize credit card: "+err.Error())
		return
	}

	resp := authorizeResponse{
		UID:      req.UID,
		Amount:   m.Amount(),
		Currency: m.Currency(),
	}

	body, err = json.Marshal(resp)
	if err != nil {
		internalError(w, "failed to marshal response")
		return
	}

	responseOK(w, body)
}
