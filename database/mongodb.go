package database

import (
	"NebuloGo/config"
	"NebuloGo/salt"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Database   *mongo.Database
}

// NewUserManager crée un nouveau UserManager
func NewUserManager(collection *mongo.Collection) *UserManager {
	return &UserManager{collection: collection, Database: collection.Database()}
}

// CreateUser crée un nouvel utilisateur
func (um *UserManager) CreateUser(loginID, password string) error {
	user := MongoUser{
		InternalID:     uuid.NewString(),
		LoginID:        loginID,
		HashedPassword: salt.HashPhrase(password),
		LastModified:   time.Now(),
	}

	_, err := um.collection.InsertOne(context.TODO(), user)
	return err
}

// UpdateLoginID met à jour l'identifiant de connexion d'un utilisateur
func (um *UserManager) UpdateLoginID(internalID, newLoginID string) error {
	// Mettre à jour l'identifiant de connexion
	_, err := um.collection.UpdateOne(context.TODO(),
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

func MongoDBInit() error {
	// Configuration du client MongoDB
	clientOptions := options.Client().ApplyURI(config.Configuration.Database.ServerURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Connexion à la base de données et à la collection "users"
	collection := client.Database(config.Configuration.Database.DatabaseName).Collection("users")

	// Créer un index unique sur le champ internal_id
	indexModelInternalID := mongo.IndexModel{
		Keys:    bson.M{"internal_id": 1}, // 1 signifie un ordre croissant
		Options: options.Index().SetUnique(true),
	}

	// Créer un index unique sur le champ login_id
	indexModelLoginID := mongo.IndexModel{
		Keys:    bson.M{"login_id": 1}, // 1 signifie un ordre croissant
		Options: options.Index().SetUnique(true),
	}

	// Créer les deux index
	_, err = collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{indexModelInternalID, indexModelLoginID})
	if err != nil {
		return err
	}

	// Créer un gestionnaire d'utilisateurs
	ApplicationUserManager = NewUserManager(collection)
	return nil
}
