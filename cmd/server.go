package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"truckx/internal/config"
	"truckx/internal/controllers"
	"truckx/internal/db"
	"truckx/internal/middlewares"
	"truckx/internal/models"
	"truckx/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Initialize Config
	cfg := config.Load()

	// Initialize Database
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalln("db.Connect() failed, err:", err.Error())
	}

	// Initialize Router
	router := initRouter(dbConn, cfg)

	// Initialize HTTP Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start HTTP Server
	go func() {
		log.Println("listening...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to listen: %s\n", err)
		}
	}()

	// Await interruption
	deathRay := make(chan os.Signal, 1)
	signal.Notify(deathRay, syscall.SIGINT, syscall.SIGTERM)
	<-deathRay
	log.Println("gracefull shutdown of server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error", err)
	}
}

func initRouter(dbConn *sqlx.DB, cfg models.Config) *gin.Engine {
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"app": "theatre booking engine"})
	})

	// un-secured routes
	api := engine.Group("/api")

	authCtrl := controllers.NewAuthCtrl(repository.NewAuthRepo(dbConn), cfg.JWTKet)
	api.POST("/login", authCtrl.Login)
	api.POST("/signup", authCtrl.Signup)

	// secured routes
	apiV1 := api.Group("/v1").Use(middlewares.Auth(cfg.JWTKet))

	movieCtrl := controllers.NewMovieCtrl(repository.NewMovieRepo(dbConn))
	apiV1.POST("/movies", movieCtrl.Create)
	apiV1.GET("/movies", movieCtrl.List)
	apiV1.GET("/movies/:id/availability", movieCtrl.GetAvailability)

	bookingCtrl := controllers.NewBookingCtrl(repository.NewMovieRepo(dbConn), repository.NewBookingRepo(dbConn))
	apiV1.POST("/bookings", bookingCtrl.Create)
	apiV1.GET("/bookings", bookingCtrl.List)
	apiV1.GET("/bookings/:id", bookingCtrl.Get)

	return engine
}
