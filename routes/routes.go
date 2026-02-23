package routes

import (
	"go-api/backend-api/controllers"
	"go-api/backend-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	// route cross origin resource sharing (CORS)
	router.Use(middlewares.CORSMiddleware())

	// route register
	router.POST("/api/register", controllers.Register)

	// route login
	router.POST("/api/login", controllers.Login)

	// route get users
	router.GET("/api/users", middlewares.AuthMiddleware(), controllers.FindUsers)

	// route get user by id
	router.GET("/api/users/:id", middlewares.AuthMiddleware(), controllers.FindUserById)

	// route post create user
	router.POST("/api/users", middlewares.AuthMiddleware(), controllers.CreateUser)

	// route Patch update user
	router.PATCH("/api/users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)

	// route delete user
	router.DELETE("/api/users/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)

	return router
}
