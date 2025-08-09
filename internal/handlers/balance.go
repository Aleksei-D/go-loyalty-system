package handlers

import (
	"encoding/json"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"net/http"
)

type BalanceHandler struct {
	bs *service.BalanceService
}

func NewBalanceHandler(bs *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{bs: bs}
}

func (b *BalanceHandler) ApiGetBalanceHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login, _ := r.Context().Value("login").(string)

		balance, err := b.bs.Get(r.Context(), login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		balanceJson, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "service/json")
		w.WriteHeader(http.StatusOK)
		w.Write(balanceJson)
	}
}
