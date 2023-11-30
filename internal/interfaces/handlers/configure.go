package handlers

import (
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/interfaces/middlewares"
	"github.com/go-chi/chi/v5"
)

func UserRouter(uh *HandlerUser, secretKey []byte) chi.Router {
	router := chi.NewRouter()

	authentication := middlewares.Authentication(secretKey)

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", uh.RegisterUser)
		r.Post("/login", uh.Login)
		r.With(authentication).Post("/orders", uh.AddOrder)
		r.With(authentication).Get("/orders", uh.GetOrders)
		r.With(authentication).Get("/balance", uh.GetBalance)
		r.With(authentication).Post("/balance/withdraw", uh.AddWithdraw)
		r.With(authentication).Get("/withdrawals", uh.GetWithdrawals)
	})

	return router
}
