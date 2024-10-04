package database

import (
	"NebuloGo/config"
	"NebuloGo/salt"
	"context"
	"errors"
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
	InternalID   primitive.ObjectID   `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	FileName     string               `bson:"file_name"`
	Owner        primitive.ObjectID   `bson:"owner_id"`
	SharedWith   []primitive.ObjectID `bson:"shared_with"`
	LastModified time.Time            `bson:"last_modified"`
	ParentID     primitive.ObjectID   `bson:"parent_id"` // Référence au dossier parent
}

type MongoFolder struct {
	InternalID   primitive.ObjectID `bson:"_id,omitempty"` // Utilisation d'ObjectID comme clé primaire
	Owner        primitive.ObjectID `bson:"owner_id"`
	FolderName   string             `bson:"folder_name"`
	LastModified time.Time          `bson:"last_modified"`
	ParentID     primitive.ObjectID `bson:"parent_id"` // Référence au dossier parent
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

type FileHierarchy struct {
	FileName   string             `bson:"file_name"`
	InternalID primitive.ObjectID `bson:"_id,omitempty"`
	Children   []*FileHierarchy   `bson:"children"` // Utilisation de pointeurs ici
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

	indexModelFiles := mongo.IndexModel{
		Keys: bson.D{
			{"user_id", 1}, // 1 pour l'index ascendant
			{"name", 1},
			{"path", 1},
		},
		Options: options.Index().SetUnique(true), // Spécifier que l'index est unique
	}
	// Créer l'index
	_, err = filesCollection.Indexes().CreateOne(context.TODO(), indexModelFiles)
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

// CreateFile crée un nouveau fichier (parent vide = à la racine de l'utilisateur)
func (fm *FileManager) CreateFile(ownerID primitive.ObjectID, filename string, parentID primitive.ObjectID) error {
	// Si parentID n'est pas spécifié (vide), alors il est à la racine
	if parentID.IsZero() {
		parentID = primitive.NilObjectID // Spécifie qu'il est à la racine
	}

	file := MongoFile{
		InternalID:   primitive.NewObjectID(), // Génération automatique de l'ObjectID
		Owner:        ownerID,
		FileName:     filename,
		SharedWith:   []primitive.ObjectID{},
		LastModified: time.Now(),
		ParentID:     parentID,
	}

	_, err := fm.collection.InsertOne(context.TODO(), file)
	return err
}

// ListFiles récupère tous les fichiers dans un dossier
func (fm *FileManager) ListFiles(parentID primitive.ObjectID) ([]MongoFile, error) {
	var files []MongoFile
	cursor, err := fm.collection.Find(context.TODO(), bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var file MongoFile
		if err := cursor.Decode(&file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// DeleteFile supprime un fichier appartenant à un utilisateur par son InternalID
func (fm *FileManager) DeleteFile(internalID primitive.ObjectID) error {
	_, err := fm.collection.DeleteOne(context.TODO(), bson.M{"_id": internalID})
	return err
}

// RenameFile renomme un fichier appartenant à un utilisateur
func (fm *FileManager) RenameFile(internalID primitive.ObjectID, newFileName string) error {
	_, err := fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"file_name": newFileName, "last_modified": time.Now()}})
	return err
}

// MoveFile déplace un fichier appartenant à un utilisateur en modifiant son id de parent
func (fm *FileManager) MoveFile(internalID primitive.ObjectID, newParentID primitive.ObjectID) error {
	_, err := fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"parent_id": newParentID, "last_modified": time.Now()}})
	return err
}

// GetFileHierarchy récupère la hiérarchie des fichiers pour un utilisateur donné par son internalID
func (fm *FileManager) GetFileHierarchy(userID primitive.ObjectID) ([]FileHierarchy, error) {
	// Récupérer tous les dossiers de l'utilisateur
	foldersCursor, err := fm.collection.Find(context.TODO(), bson.M{"owner_id": userID})
	if err != nil {
		return nil, err
	}
	defer foldersCursor.Close(context.TODO())

	var folderList []MongoFolder
	if err := foldersCursor.All(context.TODO(), &folderList); err != nil {
		return nil, err
	}

	// Récupérer tous les fichiers de l'utilisateur
	filesCursor, err := fm.collection.Find(context.TODO(), bson.M{"owner_id": userID})
	if err != nil {
		return nil, err
	}
	defer filesCursor.Close(context.TODO())

	var fileList []MongoFile
	if err := filesCursor.All(context.TODO(), &fileList); err != nil {
		return nil, err
	}

	// Créer une carte pour les dossiers pour un accès rapide
	folderMap := make(map[primitive.ObjectID]*FileHierarchy) // Utilisation de pointeurs ici
	for _, folder := range folderList {
		folderMap[folder.InternalID] = &FileHierarchy{ // Utilisation de & pour créer un pointeur
			FileName:   folder.FolderName,
			InternalID: folder.InternalID,
			Children:   []*FileHierarchy{}, // Initialise comme une liste vide
		}
	}

	// Organiser les fichiers dans leurs dossiers respectifs
	for _, file := range fileList {
		if folder, exists := folderMap[file.ParentID]; exists {
			// Ajout du fichier à la liste des enfants du dossier
			folder.Children = append(folder.Children, &FileHierarchy{ // Utilisation de & pour créer un pointeur
				FileName:   file.FileName,
				InternalID: file.InternalID,
			})
		}
	}

	// Récupérer les dossiers racines (sans parent)
	var hierarchy []FileHierarchy
	for _, folder := range folderList {
		if folder.ParentID.IsZero() { // Vérifie si c'est un dossier à la racine
			hierarchy = append(hierarchy, *folderMap[folder.InternalID]) // Déréférencer le pointeur pour ajouter la valeur
		}
	}

	return hierarchy, nil
}

// --------------------------- MÉTHODES POUR FOLDER MANAGER ---------------------------

// CreateFolder crée un nouveau dossier (parent vide = à la racine de l'utilisateur)
func (fm *FolderManager) CreateFolder(ownerID primitive.ObjectID, folderName string, parentID primitive.ObjectID) error {
	// Si parentID n'est pas spécifié (vide), alors il est à la racine
	if parentID.IsZero() {
		parentID = primitive.NilObjectID // Spécifie qu'il est à la racine
	}

	folder := MongoFolder{
		InternalID:   primitive.NewObjectID(), // Génération automatique de l'ObjectID
		Owner:        ownerID,
		FolderName:   folderName,
		LastModified: time.Now(),
		ParentID:     parentID,
	}

	_, err := fm.collection.InsertOne(context.TODO(), folder)
	return err
}

// ListFolder récupère tous les dossiers dans un dossier
func (fm *FolderManager) ListFolder(parentID primitive.ObjectID) ([]MongoFolder, error) {
	var folders []MongoFolder
	cursor, err := fm.collection.Find(context.TODO(), bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var folder MongoFolder
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// DeleteFolder supprime un dossier appartenant à un utilisateur par son InternalID
func (fm *FolderManager) DeleteFolder(internalID primitive.ObjectID) error {
	_, err := fm.collection.DeleteOne(context.TODO(), bson.M{"_id": internalID})
	return err
}

// RenameFolder renomme un dossier appartenant à un utilisateur
func (fm *FolderManager) RenameFolder(internalID primitive.ObjectID, newFolderName string) error {
	_, err := fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"folder_name": newFolderName, "last_modified": time.Now()}})
	return err
}

// MoveFolder déplace un dossier appartenant à un utilisateur en modifiant son id de parent
// MoveFolder déplace un dossier appartenant à un utilisateur en modifiant son id de parent
func (fm *FolderManager) MoveFolder(internalID primitive.ObjectID, newParentID primitive.ObjectID) error {
	// Vérifier si le nouveau parent est un ancêtre du dossier à déplacer
	isDescendant, err := fm.IsDescendant(internalID, newParentID)
	if err != nil {
		return err // Retourner une erreur si la vérification échoue
	}
	if isDescendant {
		return errors.New("un dossier ne peut pas être déplacé dans l'un de ses sous-dossiers")
	}

	_, err = fm.collection.UpdateOne(context.TODO(),
		bson.M{"_id": internalID},
		bson.M{"$set": bson.M{"parent_id": newParentID, "last_modified": time.Now()}})
	return err
}

// IsDescendant vérifie si le dossier donné est un descendant d'un autre dossier
func (fm *FolderManager) IsDescendant(folderID, potentialParentID primitive.ObjectID) (bool, error) {
	currentID := folderID

	for {
		var folder MongoFolder
		err := fm.collection.FindOne(context.TODO(), bson.M{"_id": currentID}).Decode(&folder)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return false, nil // Si le dossier n'existe pas, retourner faux
			}
			return false, err // Retourner une erreur si une autre erreur s'est produite
		}

		// Vérifier si le parent du dossier est celui que nous examinons
		if folder.ParentID == potentialParentID {
			return true, nil // Le dossier est un descendant
		}

		// Mettre à jour l'ID courant pour passer au parent
		currentID = folder.ParentID
	}
}
