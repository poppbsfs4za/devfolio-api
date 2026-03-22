package main

// @title Devfolio API
// @version 1.0
// @description Personal dev platform API
// @host localhost:8080
// @BasePath /api/v1

import (
	"log"

	_ "github.com/example/devfolio-api/docs"
	"github.com/example/devfolio-api/internal/config"
	"github.com/example/devfolio-api/internal/database"

	"github.com/example/devfolio-api/internal/delivery/http/handlers"
	"github.com/example/devfolio-api/internal/domain/entities"
	"github.com/example/devfolio-api/internal/infrastructure/persistence/repository"
	"github.com/example/devfolio-api/internal/router"
	"github.com/example/devfolio-api/internal/usecase"
	pkgAuth "github.com/example/devfolio-api/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// if cfg.App.AutoMigrate {
	// 	if err := database.AutoMigrate(db); err != nil {
	// 		log.Fatalf("failed to automigrate: %v", err)
	// 	}
	// }

	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	tagRepo := repository.NewTagRepository(db)
	postRepo := repository.NewPostRepository(db)

	if err := seedAdmin(userRepo, cfg); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiresInHours)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	tagUsecase := usecase.NewTagUsecase(tagRepo)
	postUsecase := usecase.NewPostUsecase(postRepo, tagRepo)

	h := router.Handlers{
		Health:  handlers.NewHealthHandler(),
		Auth:    handlers.NewAuthHandler(authUsecase),
		Posts:   handlers.NewPostHandler(postUsecase),
		Tags:    handlers.NewTagHandler(tagUsecase),
		Profile: handlers.NewProfileHandler(profileUsecase),
		Project: handlers.NewProjectHandler(projectUsecase),
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	router.Setup(app, h, cfg.JWT.Secret)

	log.Printf("%s running on :%s", cfg.App.Name, cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedAdmin(userRepo *repository.UserRepository, cfg *config.Config) error {
	existing, err := userRepo.GetByEmail(cfg.Seed.AdminEmail)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}
	hash, err := pkgAuth.HashPassword(cfg.Seed.AdminPassword)
	if err != nil {
		return err
	}
	return userRepo.Create(&entities.User{
		Email:        cfg.Seed.AdminEmail,
		PasswordHash: hash,
		DisplayName:  cfg.Seed.AdminDisplayName,
	})
}
