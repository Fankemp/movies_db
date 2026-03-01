package app

import (
	"context"
	"filmDb/internal/handlers"
	"filmDb/internal/repository/postgres"
	"filmDb/internal/repository/postgres/movies"
	"filmDb/pkg/modules"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("there is no .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()

	pg, err := postgres.NewStorage(ctx, dbConfig)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	movieRepo := movies.NewRepository(pg)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	r := gin.Default()

	r.POST("/movies", movieHandler.Create)
	r.GET("/movies", movieHandler.GetAllMovies)
	r.GET("/movie/search", movieHandler.Search)
	r.GET("/movie/:id", movieHandler.GetMovieById)
	r.PATCH("/movie", movieHandler.UpdateRating)
	r.DELETE("/movie/:id", movieHandler.DeleteMovieByTitle)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	err = pg.DB.Close()
	if err != nil {
		log.Fatalf("problem %v", err)
	}
	log.Println("server exited")

}

func initPostgreConfig() *modules.PostgreConfig {
	timeoutRaw := os.Getenv("DB_EXEC_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutRaw)
	if err != nil {
		timeout = 5 * time.Second
	}
	return &modules.PostgreConfig{
		HOST:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Username:    os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		SSLMode:     os.Getenv("DB_SSLMODE"),
		ExecTimeout: timeout,
	}
}
