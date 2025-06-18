package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"go-vue-auth-crud/config"
	"go-vue-auth-crud/controllers"
	"go-vue-auth-crud/routes"
)

func setupRouter() (*gin.Engine, *gorm.DB) {
	cfg := config.Load()
	db := config.ConnectDB(cfg)
	// use sqlite memory for tests perhaps, skipped for brevity
	r := gin.Default()
	routes.RegisterRoutes(r, db, cfg)
	return r, db
}

func TestRegisterEndpoint(t *testing.T) {
	r, db := setupRouter()
	defer db.Close()

	body := map[string]string{"name": "Tester", "email": "test@x.com", "password": "pass123"}
	jsonVal, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}