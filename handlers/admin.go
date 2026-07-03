package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminDashboard(c *gin.Context) {
	h.renderPage(c, http.StatusOK, "admin/dashboard/index.html", gin.H{}, PageOptions{
		Title: "Dashboard",
	})
}
