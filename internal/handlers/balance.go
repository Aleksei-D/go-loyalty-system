package handlers

import (
	"encoding/json"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
	"net/http"
)

type BalanceHandler struct {
	bs *service.BalanceService
}

func NewBalanceHandler(bs *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{bs: bs}
}

func (b *BalanceHandler) APIGetBalanceHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		balance, err := b.bs.Get(r.Context(), login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(balanceJSON)
	}
}
