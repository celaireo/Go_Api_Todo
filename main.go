package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"time"
)

// Structure d'une tâche
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Base de données en mémoire (tableau de tâches)
var tasks []Task
var nextID = 1

const filePath = "tasks.json"

// Charger les tâches depuis tasks.json
func loadTasksFromFile() {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Aucune tâche existante ou erreur de lecture.")
		return
	}

	if len(file) == 0 {
		fmt.Println("Fichier tasks.json vide.")
		return
	}

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		fmt.Println("Erreur lors du chargement des tâches :", err)
	}

	// Mettre à jour nextID pour éviter les doublons
	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}
}

// Sauvegarder les tâches dans tasks.json
func saveTasksToFile() {
	file, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement des tâches :", err)
		return
	}

	err = ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier :", err)
	}
}

func addTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
		return
	}
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	saveTasksToFile() // Sauvegarde après ajout
	c.JSON(http.StatusCreated, newTask)
}

//--------------Utilisation des Goroutines-----------------------

// Route pour simuler un traitement en arrière-plan
func processTask(c *gin.Context) { 

	c.JSON(http.StatusAccepted, gin.H{"message": "Traitement en arrière-plan lancé"})

	// Lancement de la goroutine immédiatement après l'envoi de la réponse au client
	go func() {
		fmt.Println("Début du traitement de la tâche...")
		time.Sleep(5 * time.Second) 
		fmt.Println("Traitement terminé !")
	}()
}

func main() {
	r := gin.Default()

	// Charger les tâches existantes
	loadTasksFromFile()
	fmt.Println("Tâches chargées :", tasks)

	// Récupérer toutes les tâches
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	})

	// Ajouter une tâche
	r.POST("/tasks", addTask)

	// Modifier une tâche
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

		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Title = updatedTask.Title
				tasks[i].Done = updatedTask.Done
				saveTasksToFile() // Sauvegarde après modification
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
	})

	// Supprimer une tâche
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
			return
		}

		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				saveTasksToFile() // Sauvegarde après suppression
				c.JSON(http.StatusOK, gin.H{"message": "Tâche supprimée"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
	})

	//on déclare une nouvelle route HTTP associée à processTask
	r.GET("/tasks/process", processTask)

	// Lancer le serveur
	r.Run(":8080")
}
