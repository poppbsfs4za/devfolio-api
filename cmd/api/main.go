package main

// @title Devfolio API
// @version 1.0
// @description Personal dev platform API
// @host localhost:8080
// @BasePath /api/v1

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/storage"
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
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	log.Printf("[startup] loading config...")
	cfg := config.Load()

	log.Printf("[startup] connecting to database (host=%s port=%s db=%s sslmode=%s)...",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	// Up to 3 attempts to smooth over a transient/cold-start blip (e.g. Neon
	// waking from auto-suspend). If the database is still unreachable after
	// that - for example the monthly Neon compute quota is exhausted - we do
	// NOT crash the process. Instead the API starts in a degraded state:
	// /health still reports OK (liveness), /ready(z) reports 503, and every
	// DB-dependent route returns a structured 503 via the RequireDB guard
	// middleware. This keeps Cloud Run's health checks meaningful instead of
	// the whole container failing to start.
	db, err := database.ConnectWithRetry(cfg.DB, cfg.App.Env, 3, 2*time.Second)
	dbAvailable := err == nil
	if err != nil {
		log.Printf("[startup] WARNING: database unavailable after retries, starting in degraded mode: %v", err)
	} else {
		log.Printf("[startup] database connection established")
	}

	dbStatus := database.NewStatus(db)

	if dbAvailable && cfg.App.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			// A migration failure with a live DB connection indicates a real
			// schema/config problem, not a transient quota issue - safe to
			// fail fast here as before.
			log.Fatalf("failed to automigrate: %v", err)
		}
	} else if !dbAvailable {
		log.Printf("[startup] skipping automigrate: database unavailable")
	}

	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	tagRepo := repository.NewTagRepository(db)
	postRepo := repository.NewPostRepository(db)

	if dbAvailable {
		if err := seedAdmin(userRepo, cfg); err != nil {
			log.Fatalf("failed to seed admin: %v", err)
		}
	} else {
		log.Printf("[startup] skipping admin seed: database unavailable")
	}

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiresInHours)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	tagUsecase := usecase.NewTagUsecase(tagRepo)
	postUsecase := usecase.NewPostUsecase(postRepo, tagRepo)
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	uploadUsecase := usecase.NewUploadUsecase(
		cfg.GCS.BucketName,
		cfg.GCS.PublicBaseURL,
		storageClient,
	)

	h := router.Handlers{
		Health:  handlers.NewHealthHandler(dbStatus),
		Auth:    handlers.NewAuthHandler(authUsecase),
		Posts:   handlers.NewPostHandler(postUsecase),
		Tags:    handlers.NewTagHandler(tagUsecase),
		Profile: handlers.NewProfileHandler(profileUsecase),
		Project: handlers.NewProjectHandler(projectUsecase),
		Uploads: handlers.NewUploadHandler(uploadUsecase),
	}

	app := fiber.New(fiber.Config{UnescapePath: true})
	// Recover from any panic (e.g. an unguarded nil-DB access) as a safety
	// net so a single bad request can never take down the whole process.
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     cfg.CORS.AllowMethods,
	}))
	app.Use(logger.New())
	router.Setup(app, h, cfg.JWT.Secret, dbStatus)

	log.Printf("[startup] %s starting on :%s (env=%s, db_available=%v)", cfg.App.Name, cfg.App.Port, cfg.App.Env, dbAvailable)
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
