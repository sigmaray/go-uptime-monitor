package handlers

import (
	"net/http"

	"go-uptime-monitor/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user") != nil {
		c.Redirect(http.StatusFound, "/admin/")
		return
	}
	h.renderPage(c, http.StatusOK, "admin/login/index.html", gin.H{}, PageOptions{
		Title:     "Login",
		HideNav:   true,
		BodyClass: "bg-light",
	})
}

func (h *Handler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := models.FindUserByUsername(h.DB, username)
	if err != nil {
		h.renderPage(c, http.StatusOK, "admin/login/index.html", gin.H{
			"Error": "Invalid username or password",
		}, PageOptions{
			Title:     "Login",
			HideNav:   true,
			BodyClass: "bg-light",
		})
		return
	}

	if user == nil || !models.CheckPassword(user.PasswordHash, password) {
		h.renderPage(c, http.StatusOK, "admin/login/index.html", gin.H{
			"Error": "Invalid username or password",
		}, PageOptions{
			Title:     "Login",
			HideNav:   true,
			BodyClass: "bg-light",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user", user.Username)
	session.Save()
	c.Redirect(http.StatusFound, "/admin/")
}

func (h *Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
