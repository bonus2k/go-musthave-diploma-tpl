package handlers

import "github.com/go-chi/chi/v5"

func UserRouter(uh *HandlerUser) chi.Router {
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", uh.RegisterUser)
		r.Post("/login", uh.Login)
		r.Post("/orders", uh.AddOrder)
		r.Get("/orders", uh.GetOrders)
		r.Get("/balance", uh.GetBalance)
		r.Post("/balance/withdraw", uh.AddWithdraw)
		r.Get("/withdrawals", uh.GetWithdrawals)
	})

	return router
}
