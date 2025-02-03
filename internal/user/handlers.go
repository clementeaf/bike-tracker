package user

import (
	"encoding/json"
	"net/http"

	"github.com/clementeaf/bike-tracker/internal/wallet"
	"github.com/clementeaf/bike-tracker/pkg/auth"
	httpresponse "github.com/clementeaf/bike-tracker/pkg/http"
	"github.com/clementeaf/bike-tracker/pkg/logger"
)

// POST New user
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	var input RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Error al procesar el JSON",
		})
		logger.Error("Error al decodificar JSON en /users/register", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	user, err := RegisterUser(input)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Error al registrar usuario",
		})
		logger.Error("Error al registrar usuario en /users/register", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	walletID, err := wallet.CreateDefaultWallet(user.ID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al crear la wallet",
		})
		logger.Error("Error al crear wallet en /users/register", map[string]interface{}{
			"user_id": user.ID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al generar token JWT",
		})
		logger.Error("Error al generar token JWT en /users/register", map[string]interface{}{
			"user_id": user.ID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"token":     token,
		"user_id":   user.ID.Hex(),
		"wallet_id": walletID.Hex(),
	}

	httpresponse.SendJSONResponse(w, http.StatusCreated, response)
	logger.Info("Usuario registrado y wallet creada exitosamente", map[string]interface{}{
		"user_id":   user.ID.Hex(),
		"wallet_id": walletID.Hex(),
	})
}

// POST Sign in
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Error al procesar el JSON: " + err.Error(),
		})
		return
	}

	user, err := LoginUser(creds.Email, creds.Password)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "Credenciales inválidas",
		})
		return
	}

	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al generar token JWT",
		})
		logger.Error("POST /users/login - Error al generar token JWT", map[string]interface{}{
			"user_id": user.ID.Hex(),
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, map[string]string{"token": token})
}

// GEt user info
func handleGetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	user, err := GetUserByID(userID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Usuario no encontrado",
		})
		logger.Error("GET /users/me - Usuario no encontrado", map[string]interface{}{
			"user_id": userID,
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, user)
	logger.Info("GET /users/me - Usuario encontrado", map[string]interface{}{
		"user": user,
	})
}

// DELETE user by id
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = DeleteUser(userID)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Error al eliminar usuario: " + err.Error(),
		})
		logger.Error("DELETE /users - Error al eliminar usuario", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("DELETE /users - Usuario eliminado exitosamente", map[string]interface{}{
		"user_id": userID,
	})
}

// UPDATE user
func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpresponse.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	userID, err := auth.GetAuthenticatedUserID(r)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	var input UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Error al procesar el JSON: " + err.Error(),
		})
		logger.Error("PUT /users/me/update - Error al decodificar JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	updatedUser, err := UpdateUser(userID, input)
	if err != nil {
		httpresponse.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		logger.Error("PUT /users/me/update - Error al actualizar usuario", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return
	}

	httpresponse.SendJSONResponse(w, http.StatusOK, updatedUser)
	logger.Info("PUT /users/me/update - Usuario actualizado exitosamente", map[string]interface{}{
		"user": updatedUser,
	})
}
