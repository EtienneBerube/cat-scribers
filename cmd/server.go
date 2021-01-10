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

/*
 // TODO:
 */

func RunServer() {

	c := cron.New()
	c.AddFunc("@monthly", handlers.HandleMonthlyPayments)
	c.Start()

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
		ReadTimeout:  5 * time.Second,
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

		c.Stop() // Stop Cron

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

func initRoutes(router *gin.Engine) {

	router.GET("/ping", handlers.Ping)

	router.POST("/login", handlers.Login)
	router.POST("/signup", handlers.SignUp)
	router.GET("/user/:id", handlers.GetUserByID)

	authenticated := router.Group("/") // Change authService from nil to smtg else
	{
		authenticated.Use(middleware.Auth())
		authenticated.GET("/user", handlers.GetCurrentUser)
		authenticated.PUT("/user", handlers.UpdateUser)

		authenticated.POST("/user/photo", handlers.UploadPhoto)
		authenticated.POST("/user/photos", handlers.UploadMultiplePhotos)
		authenticated.DELETE("/user/photo", handlers.DeletePhoto)
		authenticated.GET("/photo/:id", handlers.GetPhotoByID)
		authenticated.GET("/user/:id/photos", handlers.GetPhotosByOwnerID)

		authenticated.POST("/subscribe/:id", handlers.SubscribeTo)
		authenticated.DELETE("/subscribe/:id", handlers.UnsubscribeFrom)
	}
}
