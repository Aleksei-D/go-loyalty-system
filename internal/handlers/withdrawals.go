package handlers

import (
	"encoding/json"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/pkg/utils/common"
	"io"
	"net/http"
)

type WithdrawHandlers struct {
	ws *service.WithdrawalService
}

func NewWithdrawHandler(s *service.WithdrawalService) *WithdrawHandlers {
	return &WithdrawHandlers{ws: s}
}

func (wh *WithdrawHandlers) APIWithdrawHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var withdrawal models.Withdrawal
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf, &withdrawal); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if ok := common.CheckLuhnAlgorithm(withdrawal.OrderNumber); !ok {
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		withdrawal.Login = login
		err = wh.ws.Withdraw(r.Context(), &withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (wh *WithdrawHandlers) APIGetWithdrawalsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		withdrawals, err := wh.ws.GetAllByLogin(r.Context(), login)
		if len(withdrawals) == 0 {
			http.Error(w, "withdrawals not found", http.StatusNoContent)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ordersJSON, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ordersJSON)
	}
}
