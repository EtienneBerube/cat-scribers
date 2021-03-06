package cmd

import (
	"context"
	"github.com/EtienneBerube/cat-scribers/internal/handlers"
	"github.com/EtienneBerube/cat-scribers/internal/middleware"
	"github.com/EtienneBerube/cat-scribers/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// RunServer command to start the server and expose port set by env. HTTP_PORT
func RunServer() {

	c := cron.New()
	c.AddFunc("@monthly", handlers.HandleMonthlyPayments)
	c.Start()
	defer c.Stop()

	router := gin.New()
	gin.ForceConsoleColor()

	router.Use(gin.Recovery())

	// Custom Logger
	router.Use(gin.LoggerWithFormatter(middleware.WithLogging))

	address := ":" + config.Config.Port
	initRoutes(router)

	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Printf("Server is ready to handle requests at %s", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", address, err)
	}

	<-done
	log.Println("Server stopped")
}

// initRoutes will initialize the routes available to the clients
func initRoutes(router *gin.Engine) {

	router.GET("/ping", handlers.Ping)

	router.POST("/login", handlers.Login)
	router.POST("/signup", handlers.SignUp)
	router.GET("/users", handlers.GetAllUsers)
	router.GET("/user/:id", handlers.GetUserByID)

	authenticated := router.Group("/")
	{
		authenticated.Use(middleware.Auth())
		// To access the following routes, the use needs to be authenticated
		authenticated.GET("/user", handlers.GetCurrentUser)
		authenticated.PUT("/user", handlers.UpdateUser)
		authenticated.DELETE("/user", handlers.DeleteUser)

		authenticated.POST("/user/photo", handlers.UploadPhoto)
		authenticated.POST("/user/photos", handlers.UploadMultiplePhotos)
		authenticated.DELETE("/user/photo/:id", handlers.DeletePhoto)
		authenticated.GET("/photo/:id", handlers.GetPhotoByID)
		authenticated.GET("/user/:id/photos", handlers.GetPhotosByOwnerID)

		authenticated.POST("/subscribe/:id", handlers.SubscribeTo)
		authenticated.DELETE("/subscribe/:id", handlers.UnsubscribeFrom)
	}
}
