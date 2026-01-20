package store

import (
	"encoding/json"
	"fmt"
	"os"
	"server2/models"
	"server2/openai"
)

// holds all food items with their embeddings
type FoodStore struct {
	foods     []models.FoodWithEmbedding // [[name, emm], [name, emm], [name, emm]......]
	foodByID  map[string]*models.FoodWithEmbedding
	dimension int
}

// creates a new food store and loads foods from json
func NewFoodStore(dataPath string, client *openai.Client) (*FoodStore, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read foods file: %w", err)
	}

	var foods []models.Food
	if err := json.Unmarshal(data, &foods); err != nil {
		return nil, fmt.Errorf("failed to parse foods: %w", err)
	}

	store := &FoodStore{
		foods:     make([]models.FoodWithEmbedding, 0, len(foods)),
		foodByID:  make(map[string]*models.FoodWithEmbedding),
		dimension: client.GetEmbeddingDimension(),
	}

	fmt.Println("Generating embeddings for foods...")
	for i, food := range foods {
		text := food.Name + ": " + food.Description

		embedding, err := client.GetEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to get embedding for %s: %w", food.Name, err)
		}

		foodWithEmb := models.FoodWithEmbedding{
			Food:      food,
			Embedding: embedding,
		}
		store.foods = append(store.foods, foodWithEmb)
		store.foodByID[food.ID] = &store.foods[len(store.foods)-1]

		fmt.Printf("  [%d/%d] %s ✓\n", i+1, len(foods), food.Name)
	}

	fmt.Printf("Loaded %d foods with embeddings\n", len(store.foods))
	return store, nil
}

// GetAll returns all foods
func (s *FoodStore) GetAll() []models.FoodWithEmbedding {
	return s.foods
}

// food by name
func (s *FoodStore) GetByName(name string) *models.FoodWithEmbedding {
	for i := range s.foods {
		if s.foods[i].Name == name {
			return &s.foods[i]
		}
	}
	return nil
}

// FoodWithEmbedding type alias for handlers
type FoodWithEmbedding = models.FoodWithEmbedding
