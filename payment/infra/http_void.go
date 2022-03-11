package infra

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"payment-gateway/payment/domain"
)

type voidRequest struct {
	UID string `json:"uid"`
}

func (h *httpHandler) Void(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		badRequest(w, "failed to read body")
		return
	}

	req := voidRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		log.Print(err)
		badRequest(w, "failed to parse request")
		return
	}

	ctx := r.Context()

	err = h.app.Void(ctx, req.UID)
	if errors.Is(err, domain.ErrTransactionNotFound) {
		log.Print(err)
		notFound(w, "failed to capture: "+err.Error())
		return
	}

	if err != nil {
		log.Print(err)
		internalError(w, "failed to authorize credit card: "+err.Error())
		return
	}

	responseNoContent(w)
}
