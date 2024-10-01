package main

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/server"
	"NebuloGo/server/auth"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

func main() {
	config.LoadConfig()

	// Configuration du client MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongoadmin:mongopass@192.168.1.152:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Connexion à la base de données et à la collection "users"
	collection := client.Database("nebulogo").Collection("users")

	// Créer un gestionnaire d'utilisateurs
	database.ApplicationUserManager = database.NewUserManager(collection)

	if config.Configuration.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := server.NewServer()
	auth.InitJWT()
	server.ConfigureMiddlewares(app)
	server.ConfigureRoutes(app)

	err = app.Gin.SetTrustedProxies(config.Configuration.Server.TrustedProxies)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(config.Configuration.Server.Host + ":" + strconv.Itoa(config.Configuration.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
