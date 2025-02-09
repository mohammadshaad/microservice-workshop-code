package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadshaad/task-service/db"
	"github.com/mohammadshaad/task-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupTaskRoutes(router *gin.Engine) {
	tasks := router.Group("/api/tasks")
	{
		tasks.POST("", createTask)
		tasks.POST("/", createTask)
		tasks.GET("", getAllTasks)
		tasks.GET("/", getAllTasks)
		tasks.GET("/:id", getTask)
		tasks.PUT("/:id", updateTask)
		tasks.DELETE("/:id", deleteTask)
	}
}

func createTask(c *gin.Context) {
	var taskInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&taskInput); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Creating task with title: %s, description: %s, userID: %s",
		taskInput.Title, taskInput.Description, taskInput.UserID)

	task := models.NewTask(taskInput.Title, taskInput.Description, taskInput.UserID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection.InsertOne(ctx, task)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully inserted task with ID: %v", result.InsertedID)
	task.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, task)
}

func getAllTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tasks []bson.D
	cursor, err := db.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		log.Printf("Error decoding tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert BSON documents to Task structs
	var result []models.Task
	for _, doc := range tasks {
		var task models.Task
		bsonBytes, _ := bson.Marshal(doc)
		if err := bson.Unmarshal(bsonBytes, &task); err != nil {
			log.Printf("Error unmarshaling task: %v", err)
			continue
		}
		result = append(result, task)
	}

	// If no tasks found, return empty array instead of null
	if result == nil {
		result = []models.Task{}
	}

	log.Printf("Successfully retrieved %d tasks", len(result))
	c.JSON(http.StatusOK, result)
}

func getTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task models.Task
	err = db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var taskInput models.Task
	if err := c.ShouldBindJSON(&taskInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskInput.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"title":       taskInput.Title,
				"description": taskInput.Description,
				"status":      taskInput.Status,
				"updatedAt":   taskInput.UpdatedAt,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskInput)
}

func deleteTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
