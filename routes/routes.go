// routes/routes.go
package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/ImanolCE/api-rest-go/handlers"
    "github.com/ImanolCE/api-rest-go/middleware"
)

func Setup(app *fiber.App) {
    // Rutas públicas
    app.Post("/api/register", handlers.RegisterUser)
    app.Post("/api/login", handlers.LoginUser)

    // Rutas protegidas son las que requieren token JWT
    api := app.Group("/api", middleware.JWTMiddleware)
	
    // CRUD Usuarios GET/UPDATE/DELETE
    api.Get("/users", handlers.GetUsers)
    api.Get("/users/:id", handlers.GetUser)
    api.Put("/users/:id", handlers.UpdateUser)
    api.Delete("/users/:id", handlers.DeleteUser)

    // CRUD Tasks
    api.Post("/tasks", handlers.CreateTask)
    api.Get("/tasks", handlers.GetTasks)
    api.Get("/tasks/:id", handlers.GetTask)
    api.Put("/tasks/:id", handlers.UpdateTask)
    api.Delete("/tasks/:id", handlers.DeleteTask)
}
