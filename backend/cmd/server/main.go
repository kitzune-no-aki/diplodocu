package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/config"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/database"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/handlers"
	"github.com/kitzune-no-aki/diplodocu/backend/internal/utils"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, proceeding with environment variables.")
	}

	// Initialization
	utils.InitKeycloak()
	db := setupDatabase()

	// Router setup
	router := createRouter(db)
	setupRoutes(router)

	// Start server
	log.Println("Starting server on :8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupDatabase() *gorm.DB {
	dbCfg := config.LoadDBConfig()
	db, err := database.ConnectDB(dbCfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database setup completed!")
	return db
}

func createRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// CORS configuration (adjust as needed)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://diplodocu.mpech.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Inject DB into context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	return router
}

func setupRoutes(router *gin.Engine) {
	api := router.Group("/")

	// Protected routes group
	protected := api.Group("")
	protected.Use(utils.AuthMiddleware())
	{
		protected.GET("/sync-user", handlers.SyncUser)

		// Book routes
		protected.POST("/books", handlers.CreateBook)
		protected.GET("/books", handlers.ListBooks)
		protected.GET("/books/:id", handlers.GetBook)
		protected.PUT("/books/:id", handlers.UpdateBook)
		protected.DELETE("/books/:id", handlers.DeleteBook)
		// Manga routes
		protected.POST("/mangas", handlers.CreateManga)
		protected.GET("/mangas", handlers.ListMangas)
		protected.GET("/mangas/:id", handlers.GetManga)
		protected.PUT("/mangas/:id", handlers.UpdateManga)
		protected.DELETE("/mangas/:id", handlers.DeleteManga)
		// Game routes
		protected.POST("/spiel", handlers.CreateSpiel)
		protected.GET("/spiel", handlers.ListSpiele)
		protected.GET("/spiel/:id", handlers.GetSpiel)
		protected.PUT("/spiel/:id", handlers.UpdateSpiel)
		protected.DELETE("/spiel/:id", handlers.DeleteSpiel)
		// Film/serie routes
		protected.POST("/filmserie", handlers.CreateFilmserie)
		protected.GET("/filmserie", handlers.ListFilmserien)
		protected.GET("/filmserie/:id", handlers.GetFilmserie)
		protected.PUT("/filmserie/:id", handlers.UpdateFilmserie)
		protected.DELETE("/filmserie/:id", handlers.DeleteFilmserie)

		// Collection routes
		protected.POST("/sammlungen", handlers.CreateSammlung)
		protected.GET("/sammlungen", handlers.ListUserSammlungen)
		protected.GET("/sammlungen/:id", handlers.GetSammlungDetail)
		protected.DELETE("/sammlungen/:id", handlers.DeleteSammlung)

		sammlungDetail := protected.Group("/sammlung/:sammlungId")
		{
			sammlungDetail.POST("/produkte", handlers.AddProduktToSammlung)
			sammlungDetail.DELETE("/produkte/:produktId", handlers.RemoveProduktFromSammlung)

		}
	}
}
