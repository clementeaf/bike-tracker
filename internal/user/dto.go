package user

import "time"

type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type WalletInput struct {
	Email  string  `json:"email"`
	Amount float64 `json:"amount"`
}

type UserResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	WalletBalance  float64 `json:"wallet_balance"`
	LastSession    string  `json:"last_session"`
	LastBikeUsedID *string `json:"last_bike_used_id"`
}

type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func ToUserResponse(user User) UserResponse {
	var lastBikeUsedID *string
	if user.LastBikeUsedID != nil {
		id := user.LastBikeUsedID.Hex()
		lastBikeUsedID = &id
	}

	return UserResponse{
		ID:             user.ID.Hex(),
		Name:           user.Name,
		Email:          user.Email,
		WalletBalance:  user.WalletBalance,
		LastSession:    user.LastSession.Format(time.RFC3339),
		LastBikeUsedID: lastBikeUsedID,
	}
}
