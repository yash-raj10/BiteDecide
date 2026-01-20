package main

import (
	"log"
	"os"
	"server2/engine"
	"server2/handlers"
	"server2/openai"
	"server2/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// init openAi client
	openaiClient, err := openai.NewClient()
	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}

	// load food data & Generate embeddings
	dataPath := "data/foods.json"
	
	foodStore, err := store.NewFoodStore(dataPath, openaiClient)
	if err != nil {
		log.Fatalf("Failed to load foods: %v", err)
	}

	// init the components
	sessionStore := store.NewSessionStore(openaiClient.GetEmbeddingDimension())
	recommender := engine.NewRecommender(foodStore)
	handler := handlers.NewHandler(foodStore, sessionStore, recommender)

	r := gin.Default()

	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/session", handler.CreateSession)
	r.GET("/recommendation", handler.GetRecommendation)
	r.POST("/swipe", handler.Swipe)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
