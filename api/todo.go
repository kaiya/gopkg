package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaiya/gopkg/once"
	"gorm.io/gorm"
)

type todoModel struct {
	gorm.Model
	Title     string `json:"title"`
	Completed int    `json:"completed"`
}
type transformedTodo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	db  *gorm.DB
	app *gin.Engine
)

const (
	cacheControlPolicy = "max-age=40, s-maxage=60, stale-while-revalidate=3540"
)

func init() {
	//gin
	app = once.CreateApp()
	r := app.Group("/todo")
	registerRouter(r)
	//db
	db = once.OpenDB("todo", &todoModel{})
}
func registerRouter(api *gin.RouterGroup) {
	api.POST("/", CreateTodo)
	api.GET("/", FetchAllTodo)
	api.GET("/:id", FetchSingleTodo)
	api.PUT("/:id", UpdateTodo)
	api.DELETE("/:id", DeleteTodo)
}

// host:port/todo/*
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// gin.WrapF(RootHandler)
	app.ServeHTTP(w, r)
}

// 创建todo
func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

// fetchAllTodo 获取所有Todo
func FetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo
	db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	//对todos的属性做一些转换以构建更好的响应体
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	// stale-while-revalidate, vercel feature
	c.Header("Cache-Control", cacheControlPolicy)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

// fetchSingleTodo 获取单个Todo
func FetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.Header("Cache-Control", cacheControlPolicy)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo 更新单个todo
func UpdateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo 删除一个todo
func DeleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
