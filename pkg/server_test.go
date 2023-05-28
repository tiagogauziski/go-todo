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
	os.Setenv("DATABASE_URI", "todo_user:Network1@tcp(raspberrypi:31835)/todo-test?parseTime=true")

	database.ConnectDatabase(os.Getenv("DATABASE_URI"))
	database.RunMigrations()
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
		Item:      "Test item POST",
		Completed: false,
	}

	code, postTodo := postTodo(router, &todo)

	assert.Equal(t, 201, code)

	assert.Equal(t, todo.Item, postTodo.Item)
	assert.Equal(t, todo.Completed, postTodo.Completed)
}

func TestGetTodo(t *testing.T) {
	router := SetupRouter()

	todo := models.Todo{
		Item:      "Test item GET",
		Completed: false,
	}

	code, postTodo := postTodo(router, &todo)

	assert.Equal(t, 201, code)

	code, getTodo := getTodo(router, postTodo.ID)

	assert.Equal(t, 200, code)

	assert.Equal(t, todo.Item, getTodo.Item)
	assert.Equal(t, todo.Completed, getTodo.Completed)
}

func postTodo(router *gin.Engine, request *models.Todo) (int, *models.Todo) {
	w := httptest.NewRecorder()
	reqJson, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(reqJson))

	router.ServeHTTP(w, req)

	if w.Code == 201 {
		var responseTodo = models.Todo{}
		json.Unmarshal(w.Body.Bytes(), &responseTodo)

		return w.Code, &responseTodo
	}

	return w.Code, nil
}

func getTodo(router *gin.Engine, id uint) (int, *models.Todo) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos/"+strconv.FormatUint(uint64(id), 10), nil)

	router.ServeHTTP(w, req)

	if w.Code == 200 {
		var responseTodo = models.Todo{}
		json.Unmarshal(w.Body.Bytes(), &responseTodo)

		return w.Code, &responseTodo
	}

	return w.Code, nil
}
