package infra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"payment-gateway/payment/domain"
)

type voidRequest struct {
	UID string `json:"uid"`
}

func (h *httpHandler) Void(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, "failed to read body")
		return
	}

	req := voidRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		badRequest(w, "failed to parse request")
		return
	}

	ctx := r.Context()

	err = h.app.Void(ctx, req.UID)
	if errors.Is(err, domain.ErrTransactionNotFound) {
		notFound(w, "failed to capture: "+err.Error())
		return
	}

	if err != nil {
		internalError(w, "failed to authorize credit card: "+err.Error())
		return
	}

	responseNoContent(w)
}
