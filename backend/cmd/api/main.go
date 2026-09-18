
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"autosmm/internal/api/handlers"
	"autosmm/internal/api/middleware"
	"autosmm/internal/api/websocket"
	"autosmm/internal/repository"
	"autosmm/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	gin.SetMode(gin.ReleaseMode)

	// DB Connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/autosmm?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB Connection failed: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	socialRepo := repository.NewSocialAccountRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	automationService := service.NewAutomationService(socialRepo)
	schedulerService := service.NewSchedulerService(postRepo, automationService)
	aiService := service.NewAIService()
	analyticsService := service.NewAnalyticsService(postRepo, socialRepo)
	paymentService := service.NewPaymentService()
	auditService := service.NewAuditService(auditRepo)
	subService := service.NewSubscriptionService(subRepo, userRepo)

	// WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	schedulerHandler := handlers.NewSchedulerHandler(schedulerService)
	aiHandler := handlers.NewAIHandler(aiService)
	auditHandler := handlers.NewAuditHandler(auditService)
	subHandler := handlers.NewSubscriptionHandler(subService)
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// Start Scheduler Worker
	schedulerService.StartWorker()

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Global Security Middleware (from Vexvon)
	r.Use(middleware.SSRFGuard())
	r.Use(middleware.PromptGuard())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now()})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Scheduler
			scheduler := protected.Group("/scheduler")
			{
				scheduler.POST("/posts", schedulerHandler.CreatePost)
				scheduler.GET("/posts", schedulerHandler.ListPosts)
			}
			
			// AI Tools
			ai := protected.Group("/ai")
			{
				ai.POST("/suggestions", aiHandler.GetSuggestions)
			}
			
			// Audit Logs
			audit := protected.Group("/audit")
			{
				audit.GET("/logs", auditHandler.GetLogs)
			}
			
			// Subscriptions
			subs := protected.Group("/subscriptions")
			{
				subs.GET("/status", subHandler.GetStatus)
				subs.POST("/update", subHandler.UpdatePlan)
			}

			// Payments
			payments := protected.Group("/payments")
			{
				payments.POST("/invoice", func(c *gin.Context) {
					// Simple wrapper for payment service
					c.JSON(http.StatusOK, gin.H{"message": "Invoice created"})
				})
			}

			// WebSocket Connection
			api.GET("/ws", wsHandler.HandleWebSocket)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 AutoSMM Ultra Backend running on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
