package user

import (
	"context"
	"errors"
	"time"

	"github.com/clementeaf/bike-tracker/pkg/database"
	"github.com/clementeaf/bike-tracker/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// New user to databse
func RegisterUser(input RegisterUserInput) (User, error) {
	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	if err != nil {
		return User{}, err
	}
	if count > 0 {
		return User{}, errors.New("el email ya está registrado")
	}

	user := User{
		ID:             primitive.NewObjectID(),
		Name:           input.Name,
		Email:          input.Email,
		Password:       input.Password,
		WalletBalance:  0.0,
		LastSession:    time.Now(),
		LastBikeUsedID: nil,
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Login user
func LoginUser(email, password string) (User, error) {
	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return User{}, errors.New("email o contraseña incorrectos")
	} else if err != nil {
		return User{}, err
	}

	user.LastSession = time.Now()
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"last_session": user.LastSession}})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Add found to user wallet
func AddWalletBalance(input WalletInput) (User, error) {
	if input.Email == "" || input.Amount <= 0 {
		return User{}, errors.New("email válido y monto positivo son obligatorios")
	}

	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return User{}, errors.New("usuario no encontrado")
	} else if err != nil {
		return User{}, err
	}

	user.WalletBalance += input.Amount
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"wallet_balance": user.WalletBalance}})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Get user by id
func GetUserByID(userID string) (UserResponse, error) {
	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logger.Error("GetUserByID - ID inválido", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return UserResponse{}, errors.New("ID de usuario inválido")
	}

	var user User
	err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		logger.Error("GetUserByID - Usuario no encontrado", map[string]interface{}{
			"user_id": userID,
		})
		return UserResponse{}, errors.New("usuario no encontrado")
	} else if err != nil {
		logger.Error("GetUserByID - Error en la base de datos", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return UserResponse{}, err
	}

	logger.Info("GetUserByID - Usuario encontrado", map[string]interface{}{
		"user_id": userID,
		"email":   user.Email,
	})

	return ToUserResponse(user), nil
}

// Update user
func UpdateUser(userID string, input UpdateUserInput) (UserResponse, error) {
	userCollection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return UserResponse{}, errors.New("ID de usuario inválido")
	}

	// Construir actualización
	updateFields := bson.M{}
	if input.Name != nil {
		updateFields["name"] = *input.Name
	}
	if input.Email != nil {
		updateFields["email"] = *input.Email
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	if err != nil {
		return UserResponse{}, err
	}

	// Obtener el usuario actualizado
	return GetUserByID(userID)
}

// Delete user and its wallet
func DeleteUser(userID string) error {
	userCollection := database.GetCollection("users")
	walletCollection := database.GetCollection("wallets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	_, err = walletCollection.DeleteOne(ctx, bson.M{"user_id": objectID.Hex()})
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.New("error al eliminar la wallet del usuario: " + err.Error())
	}

	_, err = userCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.New("error al eliminar el usuario: " + err.Error())
	}

	return nil
}
