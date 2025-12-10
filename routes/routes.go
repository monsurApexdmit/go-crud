package routes

import (
	"go-crud/controllers"
	"github.com/go-chi/chi/v5"

)

func RegisterRoutes() *chi.Mux {
	r:= chi.NewRouter()
	r.Route("/books", func(r chi.Router) {
		r.Get("/", controllers.ListBooks)
		r.Post("/", controllers.CreateBook)
		r.Get("/{id}", controllers.GetBook)
		r.Put("/{id}", controllers.UpdateBook)
		r.Delete("/{id}", controllers.DeleteBook)
	})

	r.Route("/users",func(r chi.Router) {
		r.Get("/",controllers.ListUsers)
		r.Post("/",controllers.CreateUser)
		r.Get("/{id}",controllers.GetUser)
		r.Put("/{id}",controllers.UpdateUser)
		r.Delete("/{id}",controllers.DeleteUser)
	})

	r.Route("/authors",func(r chi.Router) {
		r.Get("/",controllers.ListAuthors)
		r.Post("/",controllers.CreateAuthor)
		r.Get("/{id}",controllers.GetAuthor)
		r.Put("/{id}",controllers.UpdateAuthor)
		r.Delete("/{id}",controllers.DeleteAuthor)
	})

	return r
}
