package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	redisClient *redis.Client
)

func main() {
	// Connexion à Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Adresse de votre instance Redis
		Password: "",               // Mot de passe Redis (si requis)
		DB:       0,                // Numéro de la base de données Redis à utiliser
	})

	// Vérifier la connexion à Redis
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Impossible de se connecter à Redis : %v", err)
	}
	log.Printf("Connexion à Redis établie. Ping : %s", pong)

	// Créer un routeur Gin
	router := gin.Default()

	// Définir les routes
	router.POST("/add", addEntry)
	router.GET("/define/:key", getEntry)
	router.GET("/remove/:key", removeEntry)
	router.GET("/list", listEntries)

	// Démarrer le serveur HTTP
	router.Run(":8080")
}

func addEntry(c *gin.Context) {
	var entry Entry
	err := c.ShouldBindJSON(&entry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données JSON invalides"})
		return
	}

	err = redisClient.Set(entry.Key, entry.Value, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de l'ajout de l'entrée dans Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrée ajoutée avec succès"})
}

func getEntry(c *gin.Context) {
	key := c.Param("key")

	value, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "La clé n'a pas été trouvée"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la récupération de l'entrée depuis Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func removeEntry(c *gin.Context) {
	key := c.Param("key")

	deleted, err := redisClient.Del(key).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la suppression de l'entrée depuis Redis"})
		return
	}

	if deleted == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "La clé n'a pas été trouvée"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrée supprimée avec succès"})
}

func listEntries(c *gin.Context) {
	keys, err := redisClient.Keys("*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la récupération des clés depuis Redis"})
		return
	}

	entries := make([]Entry, 0, len(keys))
	for _, key := range keys {
		value, _ := redisClient.Get(key).Result()
		entries = append(entries, Entry{Key: key, Value: value})
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}
