package routes

import (
	"go-crud/controllers"
	"go-crud/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	// Disable automatic redirect from /path to /path/
	// This prevents issues with POST requests containing form data
	// r.RedirectTrailingSlash = false

	r.POST("/login", controllers.Login)
	r.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)

	r.Group("/books")
	{
		r.GET("/books/", controllers.ListBooks)
		r.POST("/books/", controllers.CreateBook)
		r.GET("/books/:id", controllers.GetBook)
		r.PUT("/books/:id", controllers.UpdateBook)
		r.DELETE("/books/:id", controllers.DeleteBook)
	}

	// users := r.Group("/users", middlewares.AuthMiddleware())
	users := r.Group("/users")
	{
		users.GET("/", controllers.ListUsers)
		users.POST("/", controllers.CreateUser)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	categories := r.Group("/categories")
	{
		categories.GET("/", controllers.ListCategories)
		categories.GET("/:id", controllers.GetCategory)
		categories.POST("/", controllers.CreateCategory)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}

	attributes := r.Group("/attributes")
	{
		attributes.GET("/", controllers.ListAttributes)
		attributes.POST("/", controllers.CreateAttribute)
		attributes.GET("/:id", controllers.GetAttribute)
		attributes.PUT("/:id", controllers.UpdateAttribute)
		attributes.DELETE("/:id", controllers.DeleteAttribute)
	}

	coupons := r.Group("/coupons")
	{
		coupons.GET("/", controllers.ListCoupons)
		coupons.POST("/", controllers.CreateCoupon)
		coupons.GET("/:id", controllers.GetCoupon)
		coupons.PUT("/:id", controllers.UpdateCoupon)
		coupons.DELETE("/:id", controllers.DeleteCoupon)
	}

	locations := r.Group("/locations")
	{
		locations.GET("/", controllers.ListLocations)
		locations.POST("/", controllers.CreateLocation)
		locations.GET("/:id", controllers.GetLocation)
		locations.PUT("/:id", controllers.UpdateLocation)
		locations.DELETE("/:id", controllers.DeleteLocation)
	}
	
	productCtrl := controllers.NewProductController()

	products := r.Group("/products")
	{
		products.GET("/", productCtrl.ListProducts)
		products.POST("/", productCtrl.CreateProduct)
		products.GET("/:id", productCtrl.GetProduct)
		products.PUT("/:id", productCtrl.UpdateProduct)
		products.DELETE("/:id", productCtrl.DeleteProduct)
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
