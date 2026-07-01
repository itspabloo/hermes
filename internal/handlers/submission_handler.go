package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itspabloo/hermes/internal/models"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	DB *gorm.DB
}

func NewSubmissionHandler(db *gorm.DB) *SubmissionHandler {
	return &SubmissionHandler{DB: db}
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
	// TODO: отправка submission в очередь
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
