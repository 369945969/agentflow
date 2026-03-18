package routes

import (
	"apiServer/config"
	"apiServer/db"
	"apiServer/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterSkillRoutes(r)
	return r
}

func TestMain(m *testing.M) {
	// Change working directory to root of apiServer for test to find config.yaml
	_ = os.Chdir("..")
	
	config.LoadConfig()
	// Override for test
	config.AppConfig.Database.Path = "./data/test_agentflow.duckdb"
	
	db.InitDB()
	
	code := m.Run()
	
	db.DB.Close()
	os.Remove("./data/test_agentflow.duckdb")
	os.Exit(code)
}

func TestGetSkills(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/skills/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateSkill(t *testing.T) {
	r := setupRouter()

	skill := models.Skill{
		Name:        "Test Skill",
		Type:        "Action",
		Description: "A test skill",
		Logic:       "print('hello')",
		Version:     "1.0.0",
		Icon:        "test-icon",
		Color:       "blue",
	}
	body, _ := json.Marshal(skill)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/skills/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Skill
	err := json.Unmarshal(w.Body.Bytes(), &created)
	assert.NoError(t, err)
	assert.Equal(t, skill.Name, created.Name)
	assert.NotEmpty(t, created.ID)
}
