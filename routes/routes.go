package routes

import (
	"go-crud/controllers"
	"go-crud/middlewares"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/login", controllers.Login)
	r.POST("/logout",  middlewares.AuthMiddleware(),controllers.Logout)

	r.Group("/books")
	{
		r.GET("/books/", controllers.ListBooks)
		r.POST("/books/", controllers.CreateBook)
		r.GET("/books/:id", controllers.GetBook)
		r.PUT("/books/:id", controllers.UpdateBook)
		r.DELETE("/books/:id", controllers.DeleteBook)
	}

	r.Group("/users")
	{
		r.GET("/users/", controllers.ListUsers)
		r.POST("/users/", controllers.CreateUser)
		r.GET("/users/:id", controllers.GetUser)
		r.PUT("/users/:id", controllers.UpdateUser)
		r.DELETE("/users/:id", controllers.DeleteUser)
	}

	r.Group("/authors")
	{
		r.GET("/authors/", controllers.ListAuthors)
		r.POST("/authors/", controllers.CreateAuthor)
		r.GET("/authors/:id", controllers.GetAuthor)
		r.PUT("/authors/:id", controllers.UpdateAuthor)
		r.DELETE("/authors/:id", controllers.DeleteAuthor)
	}

	return r
}
