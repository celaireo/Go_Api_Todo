package main

import (
	"net/http"
    "strconv"

	"github.com/gin-gonic/gin"
)

// Structure d'une tâche
type Task struct {
	ID   int    `json:"id"`
	Title string `json:"title"`
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

    // Route PUT pour modifier une tâche existante
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
			return
		}

		var updatedTask Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Recherche et modification de la tâche
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Title = updatedTask.Title
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
	})

    // Lancer le serveur sur le port 8080
    r.Run(":8080")
}
