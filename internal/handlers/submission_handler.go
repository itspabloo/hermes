package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itspabloo/hermes/internal/models"
	"github.com/itspabloo/hermes/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	DB *gorm.DB
	MQ *queue.RabbitMQ
}

func NewSubmissionHandler(db *gorm.DB, mq *queue.RabbitMQ) *SubmissionHandler {
	return &SubmissionHandler{DB: db, MQ : mq}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	var submission models.Submission
	err := c.ShouldBindJSON(&submission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format: " + err.Error()})
		return
	}
	var task models.Task
	err = h.DB.First(&task, submission.TaskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate task"})
	}
	submission.Status = "Pending"
	err = h.DB.Create(&submission).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
		return
	}
	msgBody, _ := json.Marshal(map[string]uint{"submission_id": submission.ID})
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err = h.MQ.Channel.PublishWithContext(ctx,
		"",
		"submissions_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: msgBody,
		})
	if err != nil {
		log.Printf("Failed to publish a message to queue: %v", err)
	}
	c.JSON(http.StatusCreated, submission)
}

func (h *SubmissionHandler) GetSubmissiom(c *gin.Context) {
	id := c.Param("id")
	var submission models.Submission
	err := h.DB.First(&submission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
		return
	}
	c.JSON(http.StatusOK, submission)
}
