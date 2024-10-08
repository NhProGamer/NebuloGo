package database

import (
	"NebuloGo/config"
	"NebuloGo/salt"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ApplicationDataManager *DataManager

// MongoUser représente un utilisateur dans MongoDB avec ObjectID
type MongoUser struct {
	InternalID     primitive.ObjectID `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	LoginID        string             `bson:"login_id"`
	HashedPassword string             `bson:"hashed_password"`
	LastModified   time.Time          `bson:"last_modified"`
}

type ShareFile struct {
	InternalID   primitive.ObjectID   `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	FilePath     string               `bson:"file_path"`
	Owner        primitive.ObjectID   `bson:"owner_id"`
	AllowedUsers []primitive.ObjectID `bson:"allowed_users"`
	Public       bool                 `bson:"public"`
	Expiration   time.Time            `bson:"expiration"` // Champ pour la date d'expiration
}

type DataManager struct {
	Client       *mongo.Client
	UserManager  *UserManager
	ShareManager *ShareManager
}

// UserManager gère les opérations sur les utilisateurs
type UserManager struct {
	collection *mongo.Collection
}

type ShareManager struct {
	collection *mongo.Collection
}

// NewDataManager crée un nouveau DataManager
func NewDataManager(serverURL string) (*DataManager, error) {
	// Configuration du client MongoDB
	clientOptions := options.Client().ApplyURI(serverURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	usersCollection := client.Database(config.Configuration.Database.DatabaseName).Collection("users")
	sharesCollection := client.Database(config.Configuration.Database.DatabaseName).Collection("shares")

	// Créer un index unique sur le champ login_id (internal_id est géré par MongoDB)
	indexModelLoginID := mongo.IndexModel{
		Keys:    bson.M{"login_id": 1}, // 1 signifie un ordre croissant
		Options: options.Index().SetUnique(true),
	}
	// Créer l'index
	_, err = usersCollection.Indexes().CreateOne(context.TODO(), indexModelLoginID)
	if err != nil {
		return nil, err
	}

	indexModelShares := mongo.IndexModel{
		Keys: bson.M{
			"expiration": 1, // Index croissant sur le champ d'expiration
		},
		Options: options.Index().SetExpireAfterSeconds(0), // Supprimez les documents après la date d'expiration
	}

	_, err = sharesCollection.Indexes().CreateOne(context.TODO(), indexModelShares)
	if err != nil {
		return nil, err
	}

	return &DataManager{
		Client: client,
		UserManager: &UserManager{
			collection: usersCollection,
		},
		ShareManager: &ShareManager{
			collection: sharesCollection,
		},
	}, nil
}

// --------------------------- MÉTHODES POUR USER MANAGER ---------------------------

// CreateUser crée un nouvel utilisateur
func (um *UserManager) CreateUser(loginID, password string) error {
	user := MongoUser{
		InternalID:     primitive.NewObjectID(), // Génération automatique de l'ObjectID
		LoginID:        loginID,
		HashedPassword: salt.HashPhrase(password),
		LastModified:   time.Now(),
	}

	_, err := um.collection.InsertOne(context.TODO(), user)
	return err
}

// UpdateLoginID met à jour l'identifiant de connexion d'un utilisateur
func (um *UserManager) UpdateLoginID(internalID primitive.ObjectID, newLoginID string) error {
	_, err := um.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"login_id": newLoginID, "last_modified": time.Now()}})
	return err
}

// UpdatePassword met à jour le mot de passe d'un utilisateur
func (um *UserManager) UpdatePassword(internalID primitive.ObjectID, newHashedPassword string) error {
	_, err := um.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"hashed_password": newHashedPassword, "last_modified": time.Now()}})
	return err
}

// GetUserByInternalID récupère les données d'un utilisateur par son ObjectID
func (um *UserManager) GetUserByInternalID(internalID primitive.ObjectID) (*MongoUser, error) {
	var user MongoUser
	err := um.collection.FindOne(context.TODO(), bson.M{"_id": internalID}).Decode(&user)
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

// --------------------------- MÉTHODES POUR SHARE MANAGER ---------------------------

func (sm *ShareManager) CreateShare(owner primitive.ObjectID, path string, allowedUsers []primitive.ObjectID, public bool, expiration time.Time) error {
	user := ShareFile{
		InternalID:   primitive.NewObjectID(),
		FilePath:     path,
		Owner:        owner,
		AllowedUsers: allowedUsers,
		Public:       public,
		Expiration:   expiration,
	}
	_, err := sm.collection.InsertOne(context.TODO(), user)
	return err
}

func (sm *ShareManager) RemoveShare(internalID primitive.ObjectID) error {
	// Filtre basé sur l'InternalID
	filter := bson.M{"_id": internalID}

	// Supprimer le document correspondant
	_, err := sm.collection.DeleteOne(context.Background(), filter)
	return err
}

func (sm *ShareManager) GetShareFile(internalID primitive.ObjectID) (*ShareFile, error) {
	// Filtre basé sur l'InternalID
	filter := bson.M{"_id": internalID}

	// Créer une variable pour stocker le résultat
	var shareFile ShareFile

	// Rechercher le document correspondant
	err := sm.collection.FindOne(context.Background(), filter).Decode(&shareFile)
	if err != nil {
		return nil, err // Retourner l'erreur si le document n'est pas trouvé ou s'il y a un problème
	}

	return &shareFile, nil
}
