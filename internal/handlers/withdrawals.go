package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
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
		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

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

		withdrawal.Login = login
		err = wh.ws.Withdraw(r.Context(), &withdrawal)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrInvalidOrderNumber):
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			case errors.Is(err, common.ErrOrderAlreadyAdded):
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			case errors.Is(err, common.ErrPaymentInsufficient):
				http.Error(w, err.Error(), http.StatusPaymentRequired)
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
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
		if err != nil {
			if errors.Is(err, common.ErrNoContent) {
				http.Error(w, "withdrawals not found", http.StatusNoContent)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
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
