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

type MongoFile struct {
	InternalID   primitive.ObjectID `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	FileName     string             `bson:"file_name"`
	Owner        primitive.ObjectID `bson:"owner_id"`
	LastModified time.Time          `bson:"last_modified"`
	CreationDate time.Time          `bson:"creation_date"`
}

type MongoFolder struct {
	InternalID   primitive.ObjectID   `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	FolderName   string               `bson:"folder_name"`
	Owner        primitive.ObjectID   `bson:"owner_id"`
	Content      []primitive.ObjectID `bson:"content"`
	LastModified time.Time            `bson:"last_modified"`
	CreationDate time.Time            `bson:"creation_date"`
}

type DataManager struct {
	Client        *mongo.Client
	UserManager   *UserManager
	FileManager   *FileManager
	FolderManager *FolderManager
}

// UserManager gère les opérations sur les utilisateurs
type UserManager struct {
	collection *mongo.Collection
}

type FileManager struct {
	collection *mongo.Collection
}

type FolderManager struct {
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
	filesCollection := client.Database(config.Configuration.Database.DatabaseName).Collection("files")

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

	return &DataManager{
		Client: client,
		UserManager: &UserManager{
			collection: usersCollection,
		},
		FileManager: &FileManager{
			collection: filesCollection,
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

// --------------------------- MÉTHODES POUR FILE MANAGER ---------------------------

func (fm *FileManager) CreateFile(owner primitive.ObjectID, fileName string) error {
	file := MongoFile{
		InternalID:   primitive.NewObjectID(), // Génération automatique de l'ObjectID
		FileName:     fileName,
		Owner:        owner,
		LastModified: time.Now(),
	}

	_, err := fm.collection.InsertOne(context.TODO(), file)
	return err
}

func (fm *FileManager) RenameFile(fileId primitive.ObjectID, fileName string) error {
	_, err := fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": fileId},
		bson.M{"$set": bson.M{"file_name": fileName, "last_modified": time.Now()}})
	return err
}

func (fm *FileManager) RemoveFile(fileId primitive.ObjectID) error {
	// Démarrer une session pour effectuer des opérations transactionnelles
	session, err := fm.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	// Définir la logique transactionnelle
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Supprimer le fichier de la collection "files"
		_, err := fm.collection.DeleteOne(sessCtx, bson.M{"_id": fileId})
		if err != nil {
			return nil, err
		}

		// Obtenir la collection "folders" pour supprimer le fichier de tous les dossiers
		folderCollection := fm.collection.Database().Collection("folders")
		_, err = folderCollection.UpdateMany(
			sessCtx,
			bson.M{"content": fileId},
			bson.M{
				"$pull": bson.M{"content": fileId},
				"$set":  bson.M{"last_modified": time.Now()},
			},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Exécuter la transaction
	_, err = session.WithTransaction(context.TODO(), callback)
	return err
}

// --------------------------- MÉTHODES POUR FOLDER MANAGER ---------------------------

func (fm *FolderManager) CreateFolder(owner primitive.ObjectID, folderName string) error {
	file := MongoFolder{
		InternalID:   primitive.NewObjectID(), // Génération automatique de l'ObjectID
		FolderName:   folderName,
		Owner:        owner,
		Content:      []primitive.ObjectID{},
		LastModified: time.Now(),
		CreationDate: time.Now(),
	}

	_, err := fm.collection.InsertOne(context.TODO(), file)
	return err
}

func (fm *FileManager) RenameFolder(folderId primitive.ObjectID, folderName string) error {
	_, err := fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": folderId},
		bson.M{"$set": bson.M{"folder_name": folderName, "last_modified": time.Now()}})
	return err
}

func (fm *FolderManager) AddFileToFolder(folderId primitive.ObjectID, fileId primitive.ObjectID) error {
	// Mettre à jour le dossier pour ajouter l'ID du fichier à la liste de contenu
	_, err := fm.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": folderId},
		bson.M{
			"$addToSet": bson.M{"content": fileId}, // Utilise $addToSet pour éviter les doublons
			"$set":      bson.M{"last_modified": time.Now()},
		},
	)
	return err
}

func (fm *FolderManager) RemoveFileFromFolder(folderId primitive.ObjectID, fileId primitive.ObjectID) error {
	// Mettre à jour le dossier pour retirer l'ID du fichier de la liste de contenu
	_, err := fm.collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": folderId},
		bson.M{
			"$pull": bson.M{"content": fileId}, // Utilise $pull pour retirer l'ID du fichier
			"$set":  bson.M{"last_modified": time.Now()},
		},
	)
	return err
}

func (fm *FolderManager) MoveFileFromFolder(sourceFolderId primitive.ObjectID, destinationFolderId primitive.ObjectID, fileId primitive.ObjectID) error {
	err := fm.RemoveFileFromFolder(sourceFolderId, fileId)
	err = fm.AddFileToFolder(destinationFolderId, fileId)
	return err
}

func (fm *FolderManager) RemoveFolder(folderId primitive.ObjectID) error {
	// Démarrer une session pour effectuer des opérations transactionnelles
	session, err := fm.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	// Définir la logique transactionnelle
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Récupérer le dossier à supprimer
		var folder MongoFolder
		err := fm.collection.FindOne(sessCtx, bson.M{"_id": folderId}).Decode(&folder)
		if err != nil {
			return nil, err
		}

		// Supprimer récursivement le contenu du dossier (fichiers et sous-dossiers)
		for _, contentId := range folder.Content {
			// Vérifier si l'ID correspond à un fichier ou un dossier
			fileCollection := fm.collection.Database().Collection("files")
			folderCollection := fm.collection.Database().Collection("folders")

			// Chercher un fichier avec cet ID
			var file MongoFile
			fileErr := fileCollection.FindOne(sessCtx, bson.M{"_id": contentId}).Decode(&file)
			if fileErr == nil {
				// Si un fichier est trouvé, le supprimer
				_, err = fileCollection.DeleteOne(sessCtx, bson.M{"_id": contentId})
				if err != nil {
					return nil, err
				}
			} else {
				// Si aucun fichier n'est trouvé, chercher un sous-dossier avec cet ID
				var subFolder MongoFolder
				folderErr := folderCollection.FindOne(sessCtx, bson.M{"_id": contentId}).Decode(&subFolder)
				if folderErr == nil {
					// Si un sous-dossier est trouvé, le supprimer récursivement
					err = fm.RemoveFolder(contentId)
					if err != nil {
						return nil, err
					}
				} else {
					// Si l'ID ne correspond ni à un fichier ni à un dossier, retourner une erreur
					return nil, mongo.ErrNoDocuments
				}
			}
		}

		// Supprimer le dossier lui-même après avoir supprimé tout son contenu
		_, err = fm.collection.DeleteOne(sessCtx, bson.M{"_id": folderId})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Exécuter la transaction
	_, err = session.WithTransaction(context.TODO(), callback)
	return err
}
