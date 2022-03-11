package infra

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"payment-gateway/payment/domain"
)

type authorizeRequest struct {
	CCName     string `json:"cc_name"`
	CCNumber   string `json:"cc_number"`
	CCV        int    `json:"ccv"`
	ExpiryDate int    `json:"expiry_date"`
	ExpiryYear int    `json:"expiry_year"`
	Amount     int    `json:"amount"`
	Currency   string `json:"currency"`
}

type authorizeResponse struct {
	UID      string `json:"uid"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

func (h *httpHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		badRequest(w, "failed to read body")
		return
	}

	req := authorizeRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		log.Print(err)
		badRequest(w, "failed to parse request")
		return
	}

	ctx := r.Context()
	cc := domain.NewCreditCard(req.CCName, req.CCNumber, domain.NewExpiryDate(req.ExpiryDate, req.ExpiryYear), req.CCV)
	m := domain.NewMoney(req.Amount, req.Currency)

	uid, m, err := h.app.Authorize(ctx, cc, m)
	if err != nil {
		log.Print(err)
		internalError(w, "failed to authorize credit card: "+err.Error())
		return
	}

	resp := authorizeResponse{
		UID:      uid,
		Amount:   m.Amount(),
		Currency: m.Currency(),
	}

	body, err = json.Marshal(resp)
	if err != nil {
		log.Print(err)
		internalError(w, "failed to marshal response")
		return
	}

	responseOK(w, body)
}
