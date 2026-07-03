package handlers

import (
	"fmt"
	"net/http"

	"go-uptime-monitor/models"

	"github.com/gin-gonic/gin"
)

type PlaywrightSQlRequest struct {
	Query string        `json:"query" binding:"required"`
	Args  []interface{} `json:"args"`
}

type PlaywrightClearTableRequest struct {
	Table string `json:"table" binding:"required"`
}

type PlaywrightCreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) PlaywrightExecuteSQL(c *gin.Context) {
	var req PlaywrightSQlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Exec(req.Query, req.Args...).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) PlaywrightClearTable(c *gin.Context) {
	var req PlaywrightClearTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// simple way to clear table
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", req.Table)
	// fallback if sqlite doesn't support TRUNCATE, but it seems it's postgres based on config
	if err := h.DB.Exec(query).Error; err != nil {
		// if TRUNCATE fails (e.g. SQLite), try DELETE FROM
		query2 := fmt.Sprintf("DELETE FROM %s", req.Table)
		if err2 := h.DB.Exec(query2).Error; err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) PlaywrightCreateUser(c *gin.Context) {
	var req PlaywrightCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := models.CreateUserInput{
		Username:        req.Username,
		Password:        req.Password,
		ConfirmPassword: req.Password,
	}
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.FormatValidationError(err)})
		return
	}

	user, err := models.CreateUser(h.DB, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
}
