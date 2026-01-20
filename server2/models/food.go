package models

// food item with its metadata
type Food struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// food item with its embedding vector
type FoodWithEmbedding struct {
	Food
	Embedding []float64 `json:"-"`
}
