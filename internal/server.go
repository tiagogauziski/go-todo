package internal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tiagogauziski/go-todo/internal/database"
	"github.com/tiagogauziski/go-todo/internal/models"

	"github.com/gin-gonic/gin"
)

func getTodos(context *gin.Context) {
	var todos []models.Todo
	result := database.Database.Find(&todos)

	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to get all todos."})
		return
	}

	context.IndentedJSON(http.StatusOK, todos)
}

func getTodoById(id uint) (*models.Todo, error) {
	var todoModel models.Todo
	result := database.Database.First(&todoModel, id)

	if result.Error != nil {
		return nil, errors.New("invalid id")
	}

	return &todoModel, nil
}

func getTodo(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid id."})
		return
	}

	todo, err := getTodoById(uint(id))

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func addTodo(context *gin.Context) {
	var newTodo models.Todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to serialize todo."})
		return
	}

	result := database.Database.Create(&newTodo)

	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to add Todo into the database."})
	} else {
		context.IndentedJSON(http.StatusCreated, newTodo)
	}
}

func updateTodo(context *gin.Context) {
	var updateTodo models.Todo
	if err := context.BindJSON(&updateTodo); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to serialize todo."})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid id."})
	}

	todo, err := getTodoById(uint(id))
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found."})
		return
	}

	todo.Completed = updateTodo.Completed
	todo.Item = updateTodo.Item

	result := database.Database.Save(&todo)

	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update todo."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func deleteTodo(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid id."})
	}

	todo, err := getTodoById(uint(id))
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found."})
		return
	}

	database.Database.Delete(&todo)

	result := database.Database.Save(&todo)
	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to delete todo."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func toggleTodoStatus(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid id."})
		return
	}

	todo, err := getTodoById(uint(id))

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found."})
		return
	}

	todo.Completed = !todo.Completed

	result := database.Database.Save(&todo)
	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update todo."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func StartServer() {
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.POST("/todos", addTodo)

	router.Run()
}
