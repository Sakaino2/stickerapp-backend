package routes

import (
	"github.com/gin-gonic/gin"
	"jorge.com/inventari_app_backend/controllers"
	"jorge.com/inventari_app_backend/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", controllers.Signup)
	server.POST("/refresh-token", controllers.RefreshUserToken)
	server.POST("/login", controllers.Login)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/stickers", controllers.GetStickersPDF)
}
