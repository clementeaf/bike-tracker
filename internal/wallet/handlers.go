package wallet

import (
	"encoding/json"
	"net/http"

	"github.com/clementeaf/bike-tracker/pkg/auth"
	httpresponse "github.com/clementeaf/bike-tracker/pkg/http"
	"github.com/clementeaf/bike-tracker/pkg/logger"
)

// POST Add found to wallet
func HandleAddTransaction(w http.ResponseWriter, r *http.Request) {
	authenticatedUserID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("POST /wallet/transactions/add - Usuario no autenticado", nil)
		return
	}

	var input struct {
		WalletID string  `json:"wallet_id"`
		UserID   string  `json:"user_id"`
		Amount   float64 `json:"amount"`
		Type     string  `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Datos inválidos",
		})
		logger.Error("POST /wallet/transactions/add - Error al decodificar JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if input.WalletID == "" || input.UserID == "" || input.Amount <= 0 || (input.Type != "credit" && input.Type != "debit") {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Faltan datos requeridos o son inválidos",
		})
		logger.Error("POST /wallet/transactions/add - Datos faltantes o inválidos", map[string]interface{}{
			"input": input,
		})
		return
	}

	if authenticatedUserID != input.UserID {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("POST /wallet/transactions/add - Usuario autenticado no coincide", map[string]interface{}{
			"authenticated_user_id": authenticatedUserID,
			"input_user_id":         input.UserID,
		})
		return
	}

	wallet, err := GetWalletByIDAndUserID(input.WalletID, input.UserID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Wallet no encontrada o no pertenece al usuario",
		})
		logger.Error("POST /wallet/transactions/add - Wallet no encontrada o inválida", map[string]interface{}{
			"wallet_id": input.WalletID,
			"user_id":   input.UserID,
			"error":     err.Error(),
		})
		return
	}

	transaction, err := AddTransaction(wallet.ID.Hex(), input.Amount, input.Type)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		logger.Error("POST /wallet/transactions/add - Error al añadir transacción", map[string]interface{}{
			"wallet_id": input.WalletID,
			"error":     err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusCreated, transaction)
	logger.Info("POST /wallet/transactions/add - Transacción añadida exitosamente", map[string]interface{}{
		"transaction": transaction,
	})
}

// GET User wallet
func HandleGetWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("GET /wallet - Usuario no autenticado", nil)
		return
	}

	wallet, err := GetWallet(userID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Wallet no encontrada",
		})
		logger.Error("GET /wallet - Wallet no encontrada", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, wallet)
	logger.Info("GET /wallet - Wallet encontrada", map[string]interface{}{
		"user_id": userID,
		"wallet":  wallet,
	})
}

// GET User wallet transactions
func HandleGetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("GET /wallet/transactions - Usuario no autenticado", nil)
		return
	}

	transactions, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "No se encontraron transacciones",
		})
		logger.Error("GET /wallet/transactions - Sin transacciones", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, transactions)
	logger.Info("GET /wallet/transactions - Historial de transacciones obtenido", map[string]interface{}{
		"user_id":      userID,
		"transactions": len(transactions),
	})
}

// GET User wallet current balance
func HandleGetWalletBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "No autorizado",
		})
		logger.Error("GET /wallet/balance - Usuario no autenticado", nil)
		return
	}

	wallet, err := GetWallet(userID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Wallet no encontrada",
		})
		logger.Error("GET /wallet/balance - Wallet no encontrada", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, map[string]float64{"balance": wallet.Balance})
	logger.Info("GET /wallet/balance - Balance obtenido exitosamente", map[string]interface{}{
		"user_id": userID,
		"balance": wallet.Balance,
	})
}
