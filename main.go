package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/migrations/seeder"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/routes"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils/azure"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

func main() {
	rand.Seed(time.Now().Unix())

	var (
		db         *gorm.DB               = config.SetUpDatabaseConnection()
		jwtService config.JWTService      = config.NewJWTService()
		bucket     *config.SupabaseBucket = config.SetUpSupabaseBucket()

		// repositories
		userRepository          repository.UserRepository          = repository.NewUserRepository(db)
		linkShortenerRepository repository.LinkShortenerRepository = repository.NewLinkShortenerRepository(db)
		eventRepository         repository.EventRepository         = repository.NewEventRepository(db)
		pe2RSVPRepo             repository.PE2RSVPRepository       = repository.NewPE2RSVPRepository(db)
		roleRepo                repository.RoleRepository          = repository.NewRoleRepository(db)
		ticketRepository        repository.TicketRepository        = repository.NewTicketRepository(db)
		bucketRepository        repository.BucketRepository        = repository.NewSupabaseBucketRepository(bucket)

		// websocket hub
		// earlyBirdHub websocket.QueueHub = websocket.RunConnHub(eventRepository, 4, constants.MainEventEarlyBirdNoMerchID, constants.MainEventEarlyBirdWithMerchID)
		// preSaleHub   websocket.QueueHub = websocket.RunConnHub(eventRepository, 4, constants.MainEventPreSaleNoMerchID, constants.MainEventPreSaleWithMerchID)
		// normalHub    websocket.QueueHub = websocket.RunConnHub(eventRepository, 4, constants.MainEventNormalNoMerchID, constants.MainEventNormalWithMerchID)

		// services
		userService          service.UserService          = service.NewUserService(userRepository, roleRepo)
		linkShortenerService service.LinkShortenerService = service.NewLinkShortenerService(linkShortenerRepository)
		preEvent2Service     service.PreEvent2Service     = service.NewPreEvent2Service(eventRepository, pe2RSVPRepo)
		eventService         service.EventService         = service.NewEventService(eventRepository)
		mainEventService     service.MainEventService     = service.NewMainEventService(userRepository, ticketRepository, eventRepository, bucketRepository)
		storageService       service.StorageService       = service.NewStorageService(bucketRepository)

		// controllers
		userController          controller.UserController          = controller.NewUserController(userService, jwtService)
		linkShortenerController controller.LinkShortenerController = controller.NewLinkShortenerController(linkShortenerService)
		eventController         controller.EventController         = controller.NewEventController(eventService)
		preEvent2Controller     controller.PreEvent2Controller     = controller.NewPreEvent2Controller(preEvent2Service)
		mainEventController     controller.MainEventController     = controller.NewMainEventController(mainEventService)
		storageController       controller.StorageController       = controller.NewStorageController(storageService)

		// websocket handler
		// earlyBirdQueue websocket.TicketQueue = websocket.NewTicketQueue(earlyBirdHub, jwtService)
		// preSaleQueue   websocket.TicketQueue = websocket.NewTicketQueue(preSaleHub, jwtService)
		// normalQueue    websocket.TicketQueue = websocket.NewTicketQueue(normalHub, jwtService)
	)

	server := gin.Default()
	server.RedirectTrailingSlash = true

	server.Use(middleware.CORSMiddleware())

	routes.User(server, userController, jwtService)
	routes.LinkShortener(server, linkShortenerController, jwtService)
	routes.Event(server, eventController, jwtService)
	routes.PreEvent2(server, preEvent2Controller, jwtService)
	routes.MainEvent(server, mainEventController, jwtService)
	// routes.TicketQueue(server, earlyBirdQueue, preSaleQueue, normalQueue)
	routes.Storage(server, storageController, jwtService)

	// https://github.com/gin-contrib/cors
	// https://stackoverflow.com/questions/76196547/websocket-returning-403-every-time
	config := cors.DefaultConfig()
	config.AllowOrigins = constants.CORS_ALLOWED_ORIGIN
	config.AllowCredentials = true
	server.Use(cors.New(config))

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
