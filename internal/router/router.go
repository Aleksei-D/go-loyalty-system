package router

import (
	"github.com/Aleksei-D/go-loyalty-system/internal/handlers"
	"github.com/Aleksei-D/go-loyalty-system/internal/middleware"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(service *service.Service, secretKey string) chi.Router {
	userApiHandlers := handlers.NewUserHandlers(service.UserService, secretKey)
	orderApiHandlers := handlers.NewOrderHandler(service.OrderService)
	withdrawHandlers := handlers.NewWithdrawHandler(service.WithdrawalService)
	balanceHandlers := handlers.NewBalanceHandler(service.BalanceService)

	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.CompressMiddleware)
		r.Post("/register", userApiHandlers.ApiUserRegisterHandler())
		r.Post("/login", userApiHandlers.ApiUserLoginHandler())

		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("", orderApiHandlers.ApiGetOrdersHandler())
			r.Post("", orderApiHandlers.ApiAddOrdersHandler())
		})

		r.Route("/balance", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("", balanceHandlers.ApiGetBalanceHandler())
			r.Post("/withdraw", withdrawHandlers.ApiWithdrawHandler())
		})

		r.Route("/withdraws", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("", withdrawHandlers.ApiGetWithdrawalsHandler())
		})
	})
	return r
}
