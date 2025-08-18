package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
	"io"
	"net/http"
)

type OrderHandlers struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandlers {
	return &OrderHandlers{orderService: orderService}
}

func (o *OrderHandlers) APIAddOrdersHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		orderNumber := string(buf)
		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if ok := common.CheckLuhnAlgorithm(orderNumber); !ok {
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		ok, err = o.orderService.IsExist(r.Context(), orderNumber)
		if err != nil {
			http.Error(w, "adding order error", http.StatusInternalServerError)
			return
		}
		if ok {
			existOrder, err := o.orderService.GetOrderByNumber(r.Context(), orderNumber)
			if err != nil {
				http.Error(w, "adding order error", http.StatusInternalServerError)
				return
			}
			if existOrder.Login == login {
				w.WriteHeader(http.StatusOK)
				return
			}

			if existOrder.Login != login {
				http.Error(w, fmt.Sprintf("Wrong OrderNumber %s", orderNumber), http.StatusConflict)
				return
			}
		}

		_, err = o.orderService.Add(r.Context(), login, orderNumber)
		if err != nil {
			http.Error(w, "adding order error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (o *OrderHandlers) APIGetOrdersHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(common.LoginKey("login")).(string)
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		orders, err := o.orderService.GetAllByLogin(r.Context(), login)
		if len(orders) == 0 {
			http.Error(w, "Orders if not load", http.StatusNoContent)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ordersJSON, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ordersJSON)
	}
}
