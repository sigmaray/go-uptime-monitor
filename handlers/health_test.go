package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestHealthHandler(t *testing.T) {
	// Set gin to test mode so it doesn't print spam to stdout
	gin.SetMode(gin.TestMode)

	// Since we need to test DB connection, and we don't have a mock, 
	// we will just see what happens without DB or mock it loosely if possible,
	// but GORM is hard to mock. Let's just test that without valid DB it returns 503.
	db, _ := gorm.Open(postgres.Open("host=localhost user=dummy password=dummy dbname=dummy port=5432 sslmode=disable"), &gorm.Config{})
	
	h := NewHandler(db)

	r := gin.Default()
	r.GET("/health", h.Health)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// It should be 503 because the dummy DB isn't reachable
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %v", w.Code)
	}
}
