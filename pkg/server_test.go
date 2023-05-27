package pkg

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tiagogauziski/go-todo/pkg/database"
	"github.com/tiagogauziski/go-todo/pkg/models"
)

func setup() {
	os.Setenv("DATABASE_URL", "todo_user:Network1@tcp(raspberrypi:31835)/todo-test?parseTime=true")

	database.ConnectDatabase(os.Getenv("DATABASE_URL"))

	err := database.Database.AutoMigrate(&models.Todo{})

	if err != nil {
		log.Fatal("Failed to run database migrations.")
	}
}

func teardown() {
	if err := database.Database.Exec("TRUNCATE TABLE todos").Error; err != nil {
		log.Fatal("Failed to truncate `todos` table")
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestGetTodos(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPostTodo(t *testing.T) {
	router := SetupRouter()

	todo := models.Todo{
		ID:        1,
		Item:      "Test item",
		Completed: false,
	}

	w, _ := postTodo(router, &todo)

	assert.Equal(t, 201, w.Code)

	var responseTodo = models.Todo{}
	json.Unmarshal(w.Body.Bytes(), &responseTodo)

	assert.Equal(t, todo.ID, responseTodo.ID)
	assert.Equal(t, todo.Item, responseTodo.Item)
	assert.Equal(t, todo.Completed, responseTodo.Completed)
}

func TestGetTodo(t *testing.T) {
	router := SetupRouter()

	todo := models.Todo{
		ID:        2,
		Item:      "Test item",
		Completed: false,
	}

	w, _ := postTodo(router, &todo)

	assert.Equal(t, 201, w.Code)

	w, _ = getTodo(router, todo.ID)

	assert.Equal(t, 200, w.Code)

	var getTodo = models.Todo{}
	json.Unmarshal(w.Body.Bytes(), &getTodo)

	assert.Equal(t, todo.ID, getTodo.ID)
	assert.Equal(t, todo.Item, getTodo.Item)
	assert.Equal(t, todo.Completed, getTodo.Completed)
}

func postTodo(router *gin.Engine, request *models.Todo) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	json, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(json))

	router.ServeHTTP(w, req)

	return w, req
}

func getTodo(router *gin.Engine, id uint) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/"+strconv.FormatUint(uint64(id), 10), nil)

	router.ServeHTTP(w, req)

	return w, req
}
