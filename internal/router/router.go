package router

import (
	"os"

	"github.com/example/devfolio-api/internal/delivery/http/handlers"
	"github.com/example/devfolio-api/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Handlers struct {
	Health  *handlers.HealthHandler
	Auth    *handlers.AuthHandler
	Posts   *handlers.PostHandler
	Tags    *handlers.TagHandler
	Profile *handlers.ProfileHandler
	Project *handlers.ProjectHandler
}

func Setup(app *fiber.App, h Handlers, jwtSecret string) {
	if os.Getenv("APP_ENV") != "production" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}
	api := app.Group("/api/v1")

	api.Get("/health", h.Health.Health)
	api.Post("/auth/login", h.Auth.Login)

	api.Get("/profile", h.Profile.Get)
	api.Get("/posts", h.Posts.ListPublished)
	api.Get("/posts/:slug", h.Posts.GetPublishedBySlug)
	api.Get("/tags", h.Tags.List)
	api.Get("/projects", h.Project.ListFeatured)

	admin := api.Group("/admin", middleware.JWT(jwtSecret))
	admin.Post("/posts", h.Posts.Create)
	admin.Put("/posts/:id", h.Posts.Update)
	admin.Delete("/posts/:id", h.Posts.Delete)
	admin.Post("/tags", h.Tags.Create)
	admin.Put("/profile", h.Profile.Upsert)
	admin.Get("/posts", h.Posts.AdminList)
	admin.Get("/posts/:id", h.Posts.AdminGetByID)

}
