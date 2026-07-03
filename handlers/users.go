package handlers

import (
	"net/http"
	"strconv"

	"go-uptime-monitor/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UsersList(c *gin.Context) {
	var users []models.User
	h.DB.Order("created_at desc").Find(&users)

	h.renderPage(c, http.StatusOK, "admin/users/index.html", gin.H{
		"Users": users,
	}, PageOptions{Title: "Users"})
}

func (h *Handler) NewUserPage(c *gin.Context) {
	h.renderPage(c, http.StatusOK, "admin/users/create.html", gin.H{}, PageOptions{
		Title: "Create User",
	})
}

type CreateUserInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		h.renderPage(c, http.StatusBadRequest, "admin/users/create.html", gin.H{
			"Error":    "Invalid form data",
			"Username": input.Username,
		}, PageOptions{Title: "Create User"})
		return
	}

	_, err := models.CreateUser(h.DB, input.Username, input.Password)
	if err != nil {
		h.renderPage(c, http.StatusInternalServerError, "admin/users/create.html", gin.H{
			"Error":    "Failed to create user (maybe username already exists)",
			"Username": input.Username,
		}, PageOptions{Title: "Create User"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

func (h *Handler) EditUserPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	h.renderPage(c, http.StatusOK, "admin/users/edit.html", gin.H{
		"User": user,
	}, PageOptions{Title: "Edit User"})
}

type UpdateUserInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password"` // Optional, only if changing
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		h.renderPage(c, http.StatusBadRequest, "admin/users/edit.html", gin.H{
			"Error": "Invalid form data",
			"User":  user,
		}, PageOptions{Title: "Edit User"})
		return
	}

	user.Username = input.Username
	if input.Password != "" {
		hash, err := models.HashPassword(input.Password)
		if err != nil {
			h.renderPage(c, http.StatusInternalServerError, "admin/users/edit.html", gin.H{
				"Error": "Failed to hash password",
				"User":  user,
			}, PageOptions{Title: "Edit User"})
			return
		}
		user.PasswordHash = hash
	}

	if err := h.DB.Save(&user).Error; err != nil {
		h.renderPage(c, http.StatusInternalServerError, "admin/users/edit.html", gin.H{
			"Error": "Failed to update user",
			"User":  user,
		}, PageOptions{Title: "Edit User"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}
