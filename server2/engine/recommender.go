package engine

import (
	"math"
	"server2/models"
	"server2/store"
)

// swipe action weights
const (
	LeftSwipeWeight  = -0.5 // strong negative
	RightSwipeWeight = 0.2  // weak positive
	SuperSwipeWeight = 1.0  // strong positive
)

// handles food recommendation logic
type Recommender struct {
	foodStore *store.FoodStore
}

// creates a new recommender
func NewRecommender(foodStore *store.FoodStore) *Recommender {
	return &Recommender{foodStore: foodStore}
}

// returns the best unseen food for a session
func (r *Recommender) GetNextRecommendation(session *models.Session) *models.FoodWithEmbedding {
	intent := session.GetIntent()
	foods := r.foodStore.GetAll()

	var bestFood *models.FoodWithEmbedding
	bestScore := math.Inf(-1)  // start at lowest possible score

	isNeutral := IsZeroVector(intent)  // check if user has no preferences yet

	for i := range foods {
		food := &foods[i]

		if session.HasSeen(food.ID) {
			continue
		}

		var score float64
		if isNeutral {
			score = float64(len(foods) - i)
		} else {
			score = CosineSimilarity(intent, food.Embedding)
		}

		if score > bestScore {
			bestScore = score
			bestFood = food
		}
	}

	return bestFood
}

// computes cosine similarity between two vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}


//  updates the session intent based on swipe action
func (r *Recommender) UpdateIntent(session *models.Session, food *models.FoodWithEmbedding, action string) {
	intent := session.GetIntent()

	var weight float64
	switch action {
	case "left":
		weight = LeftSwipeWeight
	case "right":
		weight = RightSwipeWeight
	case "super":
		weight = SuperSwipeWeight
	default:
		return
	}

	newIntent := AddVectors(intent, ScaleVector(food.Embedding, weight))
	newIntent = NormalizeVector(newIntent)
	session.UpdateIntent(newIntent)
}


// adds two vectors element-wise
func AddVectors(a, b []float64) []float64 {
	if len(a) != len(b) {
		return a
	}
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

// multiplies a vector by a scalar (make overall vector)
func ScaleVector(v []float64, scalar float64) []float64 {
	result := make([]float64, len(v))
	for i := range v {
		result[i] = v[i] * scalar
	}
	return result
}

// normalizes a vector to unit length
func NormalizeVector(v []float64) []float64 {
	var norm float64
	for _, val := range v {
		norm += val * val
	}
	norm = math.Sqrt(norm)

	if norm == 0 {
		return v
	}

	result := make([]float64, len(v))
	for i := range v {
		result[i] = v[i] / norm
	}
	return result
}

// checks if a vector is all zeros
func IsZeroVector(v []float64) bool {
	for _, val := range v {
		if val != 0 {
			return false
		}
	}
	return true
}
