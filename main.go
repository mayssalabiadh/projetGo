// @title API Utilisateurs et Tâches
// @version 1.0
// @description Cette API permet de gérer les utilisateurs et leurs tâches
// @contact.name Développeur API
// @contact.email mayssa@gmail.com
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Saisir le token JWT comme suit : Bearer <token>

package main

import (
	"projet1/database"
	"projet1/middleware"
	"projet1/models"
	"projet1/routes"

	_ "projet1/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	godotenv.Load()
	database.Connect()

	database.DB.AutoMigrate(&models.User{}, &models.Task{}, &models.File{})

	r := gin.Default()

	//L'application globale du middleware CORS
	r.Use(middleware.CORSMiddleware())

	routes.SetupRouter(r)
	r.Static("/files", "./upload")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
