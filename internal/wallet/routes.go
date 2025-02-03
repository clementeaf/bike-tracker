package wallet

import (
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/wallet/transactions/add", middleware.AuthMiddleware(http.HandlerFunc(HandleAddTransaction)))
	mux.Handle("/wallet", middleware.AuthMiddleware(http.HandlerFunc(HandleGetWallet)))
	mux.Handle("/wallet/transactions", middleware.AuthMiddleware(http.HandlerFunc(HandleGetTransactionHistory)))
	mux.Handle("/wallet/balance", middleware.AuthMiddleware(http.HandlerFunc(HandleGetWalletBalance)))
}
