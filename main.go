package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!DOCTYPE html><html><head><title>Hello</title></head><body><h1>Hello, World!</h1></body></html>`))
	})

	r.Run(":8080")
}
