package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ApplicationUserManager *UserManager

// User struct représente un utilisateur dans MongoDB
type MongoUser struct {
	InternalID     string    `bson:"internal_id"`
	LoginID        string    `bson:"login_id"`
	HashedPassword string    `bson:"hashed_password"`
	LastModified   time.Time `bson:"last_modified"`
}

// UserManager gère les opérations sur les utilisateurs
type UserManager struct {
	collection *mongo.Collection
}

// NewUserManager crée un nouveau UserManager
func NewUserManager(collection *mongo.Collection) *UserManager {
	return &UserManager{collection: collection}
}

// CreateUser crée un nouvel utilisateur
func (um *UserManager) CreateUser(internalID, loginID, hashedPassword string) error {
	user := MongoUser{
		InternalID:     internalID,
		LoginID:        loginID,
		HashedPassword: hashedPassword,
		LastModified:   time.Now(),
	}

	_, err := um.collection.InsertOne(context.TODO(), user)
	return err
}

// UpdateLoginID met à jour l'identifiant de connexion d'un utilisateur
func (um *UserManager) UpdateLoginID(internalID, newLoginID string) error {
	// Vérifier si le nouveau login_id est déjà utilisé
	count, err := um.collection.CountDocuments(context.TODO(), bson.M{"login_id": newLoginID})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("l'identifiant de connexion '%s' est déjà utilisé", newLoginID)
	}

	// Mettre à jour l'identifiant de connexion
	_, err = um.collection.UpdateOne(context.TODO(),
		bson.M{"internal_id": internalID},
		bson.M{"$set": bson.M{"login_id": newLoginID, "last_modified": time.Now()}})
	return err
}

// UpdatePassword met à jour le mot de passe d'un utilisateur
func (um *UserManager) UpdatePassword(internalID, newHashedPassword string) error {
	_, err := um.collection.UpdateOne(context.TODO(),
		bson.M{"internal_id": internalID},
		bson.M{"$set": bson.M{"hashed_password": newHashedPassword, "last_modified": time.Now()}})
	return err
}

// GetUserByInternalID récupère les données d'un utilisateur par son ID interne
func (um *UserManager) GetUserByInternalID(internalID string) (*MongoUser, error) {
	var user MongoUser
	err := um.collection.FindOne(context.TODO(), bson.M{"internal_id": internalID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByLoginID récupère les données d'un utilisateur par son identifiant de connexion
func (um *UserManager) GetUserByLoginID(loginID string) (*MongoUser, error) {
	var user MongoUser
	err := um.collection.FindOne(context.TODO(), bson.M{"login_id": loginID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
