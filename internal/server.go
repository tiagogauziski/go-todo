package internal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tiagogauziski/go-todo/internal/database"
	"github.com/tiagogauziski/go-todo/internal/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(context *gin.Context) {
	var todos []models.Todo
	result := database.Database.Find(&todos)

	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to get all todos."})
		return
	}

	context.IndentedJSON(http.StatusOK, todos)
}

func GetTodo(context *gin.Context, todo *models.Todo) {
	context.IndentedJSON(http.StatusOK, todo)
}

func AddTodo(context *gin.Context) {
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

func UpdateTodo(context *gin.Context, todo *models.Todo) {
	var updateTodo models.Todo
	if err := context.BindJSON(&updateTodo); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to serialize todo."})
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

func DeleteTodo(context *gin.Context, todo *models.Todo) {
	database.Database.Delete(&todo)

	result := database.Database.Save(&todo)
	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to delete todo."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func ToggleTodoStatus(context *gin.Context, todo *models.Todo) {
	todo.Completed = !todo.Completed

	result := database.Database.Save(&todo)
	if result.Error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update todo."})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func getTodoWithValidationHandler(handler func(*gin.Context, *models.Todo)) gin.HandlerFunc {
	return func(context *gin.Context) {
		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid id."})
			return
		}

		var todo *models.Todo
		todo, err = getTodoById(uint(id))

		if err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found."})
			return
		}

		handler(context, todo)
	}
}

func getTodoById(id uint) (*models.Todo, error) {
	var todoModel models.Todo
	result := database.Database.First(&todoModel, id)

	if result.Error != nil {
		return nil, errors.New("invalid id")
	}

	return &todoModel, nil
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/todos", GetTodos)
	r.GET("/todos/:id", getTodoWithValidationHandler(GetTodo))
	r.PATCH("/todos/:id", getTodoWithValidationHandler(ToggleTodoStatus))
	r.PUT("/todos/:id", getTodoWithValidationHandler(UpdateTodo))
	r.DELETE("/todos/:id", getTodoWithValidationHandler(DeleteTodo))
	r.POST("/todos", AddTodo)

	return r
}
