package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Structure d'une tâche
type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Base de données en mémoire (tableau de tâches)
var tasks []Task
var nextID = 1

func main() {
	r := gin.Default()

	// Route GET pour récupérer toutes les tâches
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	})

	// Route POST pour ajouter une nouvelle tâche
	r.POST("/tasks", func(c *gin.Context) {
		var newTask Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTask.ID = nextID
		nextID++
		tasks = append(tasks, newTask)
		c.JSON(http.StatusCreated, newTask)
	})

	// Lancer le serveur sur le port 8080
	r.Run(":8080")
}
