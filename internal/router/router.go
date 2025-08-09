package router

import (
	"github.com/Aleksei-D/go-loyalty-system/internal/handlers"
	"github.com/Aleksei-D/go-loyalty-system/internal/middleware"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(service *service.Service, secretKey string) chi.Router {
	userAPIHandlers := handlers.NewUserHandlers(service.UserService, secretKey)
	orderAPIHandlers := handlers.NewOrderHandler(service.OrderService)
	withdrawHandlers := handlers.NewWithdrawHandler(service.WithdrawalService)
	balanceHandlers := handlers.NewBalanceHandler(service.BalanceService)

	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.CompressMiddleware)
		r.Post("/register", userAPIHandlers.APIUserRegisterHandler())
		r.Post("/login", userAPIHandlers.APIUserLoginHandler())

		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("/", orderAPIHandlers.APIGetOrdersHandler())
			r.Post("/", orderAPIHandlers.APIAddOrdersHandler())
		})

		r.Route("/balance", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("/", balanceHandlers.APIGetBalanceHandler())
			r.Post("/withdraw", withdrawHandlers.APIWithdrawHandler())
		})

		r.Route("/withdraws", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(secretKey))
			r.Get("/", withdrawHandlers.APIGetWithdrawalsHandler())
		})
	})
	return r
}
