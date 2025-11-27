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

	return r
}
