package handlers

import (
	"net/http"
	"errors"

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

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	err := h.DB.Preload("TestCases").First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	err := h.DB.Preload("TestCases").First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}
	err = h.DB.Where("task_id = ?", id).Delete(&models.TestCase{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test cases"})
		return
	}
	err = h.DB.Delete(&task).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
