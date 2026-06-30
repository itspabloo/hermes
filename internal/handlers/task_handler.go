package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itspabloo/hermes/internal/models"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format: " + err.Error()})
		return
	}
	err = h.DB.Create(&task).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to database"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	err := h.DB.Preload("TestCases").Find((&tasks)).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
