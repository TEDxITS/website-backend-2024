package main

import (
	"log"
	"os"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/migrations/seeder"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/routes"
	"github.com/TEDxITS/website-backend-2024/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	var (
		db         *gorm.DB           = config.SetUpDatabaseConnection()
		jwtService service.JWTService = service.NewJWTService()

		// repositories
		userRepository          repository.UserRepository          = repository.NewUserRepository(db)
		linkShortenerRepository repository.LinkShortenerRepository = repository.NewLinkShortenerRepository(db)

		// services
		userService          service.UserService          = service.NewUserService(userRepository)
		linkShortenerService service.LinkShortenerService = service.NewLinkShortenerService(linkShortenerRepository)

		// controllers
		userController          controller.UserController          = controller.NewUserController(userService, jwtService)
		linkShortenerController controller.LinkShortenerController = controller.NewLinkShortenerController(linkShortenerService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.User(server, userController, jwtService)
	routes.LinkShortener(server, linkShortenerController, jwtService)

	// database seeding, update existing data or create if not found
	if err := seeder.RunSeeders(db); err != nil {
		log.Fatalf("error migration seeder: %v", err)
	}

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	server.RedirectTrailingSlash = true
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
