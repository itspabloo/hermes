package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itspabloo/hermes/internal/handlers"
	"github.com/itspabloo/hermes/internal/repository"
)

func main() {
	db := repository.InitDB()

	taskHandler := handlers.NewTaskHandler(db)
	
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
	}
	log.Println("Starting Hermes API on port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
