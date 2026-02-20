package routes

import (
	"go-api/backend-api/controllers"
	"go-api/backend-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	// route register
	router.POST("/api/register", controllers.Register)

	// route login
	router.POST("/api/login", controllers.Login)

	// route get users
	router.GET("/api/users", middlewares.AuthMiddleware(), controllers.FindUsers)

	return router
}
