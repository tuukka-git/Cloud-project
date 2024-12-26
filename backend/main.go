package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	DB *gorm.DB
}

// Player represents data about a record player.
type Player struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" binding:"required" gorm:"not null"` // Name is required both in request and DB
}

// getPlayers responds with the list of all players from the database as JSON.
func (app *App) getPlayers(c *gin.Context) {
	var players []Player
	if result := app.DB.Find(&players); result.Error != nil {
		log.Printf("Error querying players from database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve players"})
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

// postPlayer adds a new player to the database.
func (app *App) postPlayer(c *gin.Context) {
	var newPlayer Player

	// Bind the received JSON to newPlayer with validation.
	if err := c.ShouldBindJSON(&newPlayer); err != nil {
		log.Printf("Invalid input received for new player: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Insert the new player into the database.
	if result := app.DB.Create(&newPlayer); result.Error != nil {
		log.Printf("Failed to add player '%s' to database: %v", newPlayer.Name, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add player"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

// deletePlayer removes a player from the database based on the provided ID.
func (app *App) deletePlayer(c *gin.Context) {
	idParam := c.Param("id")
	var player Player

	// Find the player by ID
	if result := app.DB.First(&player, idParam); result.Error != nil {
		log.Printf("Player not found with ID '%s': %v", idParam, result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	// Delete the player
	if result := app.DB.Delete(&player); result.Error != nil {
		log.Printf("Failed to delete player with ID '%s': %v", idParam, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player deleted"})
}

func setupRouter(app *App) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	if env := os.Getenv("ENV"); env == "DEV" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // You can customize this to restrict origins
			AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	router.GET("/players", app.getPlayers)
	router.POST("/player", app.postPlayer)
	router.DELETE("/player/:id", app.deletePlayer)
	return router
}

func createApp() (*App, error) {
	// Get the database connection URL from environment variables
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Connect to the database using GORM
	const maxRetries = 5
	var db *gorm.DB
	var err error
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			fmt.Printf("Attempt %d: Failed to connect to the database. Error: %v\n", i, err)
			time.Sleep(2 * time.Second) // wait before the next retry
		} else {
			fmt.Println("Connected to the database successfully!")
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Player{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &App{DB: db}, nil
}

func main() {
	// Start the HTTP server
	app, err := createApp()
	if err != nil {
		log.Fatalf("Error creating application: %v", err)
	}

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "0.0.0.0:8080" // Default address if not set
	}

	router := setupRouter(app)
	log.Printf("Starting server on http://%s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
