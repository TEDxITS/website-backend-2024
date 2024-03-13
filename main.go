package main

import (
	"log"
	"os"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/migrations/seeder"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/routes"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils/azure"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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
		return
	}

	/*
		Deployed on Azure App Service with .NET Stack.
		The workflow will failed to deploy on updates
		because the server is already running and it
		wont lets us replace it. Normally in .NET apps
		on azure, it will create a file called "app_offline.htm"
		and the ASP .NET will notice it the file is created
		and stop the application. This replicate said behavior.
	*/
	go azure.StopOnNewDeployment()

	port := os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = "8888"
	}

	if os.Getenv("ENV") == constants.ENUM_RUN_PRODUCTION {
		constants.BASE_URL = os.Getenv("BASE_URL")
	} else {
		constants.BASE_URL = "http://localhost:" + port
	}

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
