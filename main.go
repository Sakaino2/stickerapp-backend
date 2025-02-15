package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"jorge.com/inventari_app_backend/initializers"
	"jorge.com/inventari_app_backend/middlewares"
	"jorge.com/inventari_app_backend/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.GetDriveService()
}

func main() {
	server := gin.Default()

	server.Use(middlewares.CORSMiddleware())
	routes.RegisterRoutes(server)

	server.Run(os.Getenv("PORT")) // localhost:3000
}
