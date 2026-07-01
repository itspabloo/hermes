package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itspabloo/hermes/internal/handlers"
	"github.com/itspabloo/hermes/internal/queue"
	"github.com/itspabloo/hermes/internal/repository"
)

func main() {
	db := repository.InitDB()
	mq := queue.InitRabbitMQ()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	taskHandler := handlers.NewTaskHandler(db)
	submissionHandler := handlers.NewSubmissionHandler(db, mq)
	
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Hermes API Gateway",
			"status": "up and running",
		})
	})
	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", taskHandler.CreateTask)
		v1.GET("/tasks", taskHandler.GetTasks)
		v1.GET("/tasks/:id", taskHandler.GetTask)
		v1.DELETE("/tasks/:id", taskHandler.DeleteTask)
		v1.PUT("/tasks/:id", taskHandler.UpdateTask)

		v1.POST("/submissions", submissionHandler.CreateSubmission)
		v1.GET("/submissions/:id", submissionHandler.GetSubmissiom)
	}
	log.Printf("Starting Hermes API on port %s...", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
