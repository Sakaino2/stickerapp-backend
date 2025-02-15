package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4321", "http://localhost:8585", "http://localhost:3535", "https://inventariweb.oinc-games.com"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "Client"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowCredentials = true

	return cors.New(config)
}
