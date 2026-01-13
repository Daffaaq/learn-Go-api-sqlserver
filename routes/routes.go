// ---------------- routes/routes.go ----------------
package routes

import (
	"database/sql"
	"go-sqlserver-api/handlers"
	"go-sqlserver-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(db *sql.DB, validate *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	userRepo := &repositories.UserRepository{DB: db}
	userHandler := &handlers.UserHandler{
		Repo:      userRepo,
		Validator: validate,
	}

	r.Get("/users", userHandler.GetUsers)
	r.Post("/users", userHandler.CreateUser)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	return r
}