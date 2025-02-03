package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/clementeaf/bike-tracker/pkg/database"
	"github.com/clementeaf/bike-tracker/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create default wallet
func CreateDefaultWallet(userID primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("wallets")

	walletID := primitive.NewObjectID()
	wallet := Wallet{
		ID:          walletID,
		UserID:      userID,
		Balance:     0.0,
		LastUpdated: time.Now(),
	}

	_, err := collection.InsertOne(ctx, wallet)
	if err != nil {
		logger.Error("Error al insertar wallet en MongoDB", map[string]interface{}{
			"user_id":   userID.Hex(),
			"wallet_id": walletID.Hex(),
			"error":     err.Error(),
		})
		return primitive.NilObjectID, errors.New("falló la creación de la wallet")
	}

	logger.Info("Wallet creada exitosamente", map[string]interface{}{
		"user_id":   userID.Hex(),
		"wallet_id": walletID.Hex(),
	})
	return walletID, nil
}

// POST Found to wallet
func AddTransaction(walletID string, amount float64, transactionType string) (*Transaction, error) {
	walletCollection := database.GetCollection("wallets")
	transactionCollection := database.GetCollection("transactions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	walletObjectID, err := primitive.ObjectIDFromHex(walletID)
	if err != nil {
		return nil, errors.New("wallet ID inválido")
	}

	var wallet Wallet
	err = walletCollection.FindOne(ctx, bson.M{"_id": walletObjectID}).Decode(&wallet)
	if err != nil {
		return nil, errors.New("wallet no encontrada")
	}

	transaction := &Transaction{
		ID:        primitive.NewObjectID(),
		UserID:    wallet.UserID,
		WalletID:  walletObjectID,
		Amount:    amount,
		Type:      transactionType,
		Timestamp: time.Now(),
	}

	_, err = transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	if transactionType == "credit" {
		wallet.Balance += amount
	} else if transactionType == "debit" {
		if wallet.Balance < amount {
			return nil, errors.New("saldo insuficiente")
		}
		wallet.Balance -= amount
	} else {
		return nil, errors.New("tipo de transacción inválido")
	}

	_, err = walletCollection.UpdateOne(ctx, bson.M{"_id": walletObjectID}, bson.M{"$set": bson.M{"balance": wallet.Balance, "last_updated": time.Now()}})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GET wallet
func GetWallet(userID string) (*Wallet, error) {
	walletCollection := database.GetCollection("wallets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	var wallet Wallet
	err = walletCollection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&wallet)
	if err != nil {
		return nil, errors.New("wallet no encontrada")
	}

	return &wallet, nil
}

// GET Wallet transactions history
func GetTransactionHistory(userID string) ([]Transaction, error) {
	transactionCollection := database.GetCollection("transactions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	cursor, err := transactionCollection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// POST deduct ride fee from wallet
func DeductRideFee(userID string) error {
	const rideCost = 5.00

	walletCollection := database.GetCollection("wallets")
	transactionCollection := database.GetCollection("transactions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	var wallet Wallet
	err = walletCollection.FindOne(ctx, bson.M{"user_id": userObjectID}).Decode(&wallet)
	if err != nil {
		return errors.New("wallet no encontrada")
	}

	if wallet.Balance < rideCost {
		return errors.New("saldo insuficiente en la wallet")
	}

	_, err = walletCollection.UpdateOne(ctx, bson.M{"user_id": userObjectID}, bson.M{
		"$inc": bson.M{"balance": -rideCost},
		"$set": bson.M{"last_updated": time.Now()},
	})
	if err != nil {
		return errors.New("error al descontar saldo en la wallet")
	}

	transaction := Transaction{
		ID:        primitive.NewObjectID(),
		UserID:    userObjectID,
		Amount:    -rideCost,
		Type:      "debit",
		Timestamp: time.Now(),
	}

	_, err = transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return errors.New("error al registrar la transacción en la wallet")
	}

	return nil
}

// GET Wallet by: Wallet ID and User ID
func GetWalletByIDAndUserID(walletID, userID string) (*Wallet, error) {
	walletCollection := database.GetCollection("wallets")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	walletObjectID, err := primitive.ObjectIDFromHex(walletID)
	if err != nil {
		return nil, errors.New("wallet ID inválido")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("user ID inválido")
	}

	var wallet Wallet
	err = walletCollection.FindOne(ctx, bson.M{"_id": walletObjectID, "user_id": userObjectID}).Decode(&wallet)
	if err != nil {
		return nil, errors.New("wallet no encontrada o no pertenece al usuario")
	}

	return &wallet, nil
}
