package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/example/proofcraft/pkg/db"
	"github.com/example/proofcraft/pkg/models"
)

func setupRouter(t *testing.T) *gin.Engine {
	path := "test_ctrl.db"
	os.Remove(path)
	dbConn, err := db.Init(path)
	if err != nil {
		t.Fatalf("db init: %v", err)
	}
	dbConn.AutoMigrate(&models.Key{}, &models.Proof{})
	ctrl := NewController(dbConn)
	r := gin.New()
	r.POST("/keys/generate", ctrl.GenerateKey)
	return r
}

func TestGenerateKeyEndpoint(t *testing.T) {
	r := setupRouter(t)

	req := `{"label":"itest"}`
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("POST", "/keys/generate", strings.NewReader(req)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["private_key"] == nil {
		t.Fatalf("expected private_key in response")
	}
}
