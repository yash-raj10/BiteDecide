package engine

import (
	"math"
	"server2/models"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1, 0, 0},
			b:        []float64{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			a:        []float64{1, 1, 0},
			b:        []float64{1, 0, 0},
			expected: 1 / math.Sqrt(2),
		},
		{
			name:     "zero vector a",
			a:        []float64{0, 0, 0},
			b:        []float64{1, 0, 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("CosineSimilarity() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddVectors(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	expected := []float64{5, 7, 9}

	result := AddVectors(a, b)

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("AddVectors()[%d] = %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestScaleVector(t *testing.T) {
	v := []float64{2, 4, 6}
	scalar := 0.5
	expected := []float64{1, 2, 3}

	result := ScaleVector(v, scalar)

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("ScaleVector()[%d] = %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestNormalizeVector(t *testing.T) {
	v := []float64{3, 4, 0}
	result := NormalizeVector(v)

	var mag float64
	for _, val := range result {
		mag += val * val
	}
	mag = math.Sqrt(mag)

	if math.Abs(mag-1.0) > 0.0001 {
		t.Errorf("NormalizeVector() magnitude = %v, want 1.0", mag)
	}
}

func TestNormalizeZeroVector(t *testing.T) {
	v := []float64{0, 0, 0}
	result := NormalizeVector(v)

	for i, val := range result {
		if val != 0 {
			t.Errorf("NormalizeVector(zero)[%d] = %v, want 0", i, val)
		}
	}
}

func TestIntentUpdateLeft(t *testing.T) {
	session := models.NewSession("test", 3)
	session.UpdateIntent([]float64{0.5, 0.5, 0})

	food := &models.FoodWithEmbedding{
		Food:      models.Food{ID: "1", Name: "Test"},
		Embedding: []float64{1, 0, 0},
	}

	r := &Recommender{}

	initialIntent := session.GetIntent()
	r.UpdateIntent(session, food, "left")
	newIntent := session.GetIntent()

	initialSim := CosineSimilarity(initialIntent, food.Embedding)
	newSim := CosineSimilarity(newIntent, food.Embedding)

	if newSim >= initialSim {
		t.Errorf("Left swipe should decrease similarity: was %v, now %v", initialSim, newSim)
	}
}

func TestIntentUpdateRight(t *testing.T) {
	session := models.NewSession("test", 3)
	session.UpdateIntent([]float64{0, 1, 0})

	food := &models.FoodWithEmbedding{
		Food:      models.Food{ID: "1", Name: "Test"},
		Embedding: []float64{1, 0, 0},
	}

	r := &Recommender{}

	r.UpdateIntent(session, food, "right")
	newIntent := session.GetIntent()

	sim := CosineSimilarity(newIntent, food.Embedding)
	if sim <= 0 {
		t.Errorf("Right swipe should increase similarity toward food, got %v", sim)
	}
}

func TestIntentUpdateSuper(t *testing.T) {
	session := models.NewSession("test", 3)
	session.UpdateIntent([]float64{0, 1, 0})

	food := &models.FoodWithEmbedding{
		Food:      models.Food{ID: "1", Name: "Test"},
		Embedding: []float64{1, 0, 0},
	}

	r := &Recommender{}

	r.UpdateIntent(session, food, "super")
	newIntent := session.GetIntent()

	sim := CosineSimilarity(newIntent, food.Embedding)
	if sim < 0.5 {
		t.Errorf("Super swipe should strongly reinforce food direction, got similarity %v", sim)
	}
}

func TestSessionDoesNotRepeatFoods(t *testing.T) {
	session := models.NewSession("test", 3)

	session.MarkSeen("1")
	session.MarkSeen("2")

	if !session.HasSeen("1") {
		t.Error("Session should have seen food 1")
	}
	if !session.HasSeen("2") {
		t.Error("Session should have seen food 2")
	}
	if session.HasSeen("3") {
		t.Error("Session should not have seen food 3")
	}
}

func TestSessionCompletion(t *testing.T) {
	session := models.NewSession("test", 3)

	if session.IsCompleted() {
		t.Error("New session should not be completed")
	}

	session.Complete("Pizza")

	if !session.IsCompleted() {
		t.Error("Session should be completed after Complete()")
	}
}

func TestZeroVectorDetection(t *testing.T) {
	zero := []float64{0, 0, 0}
	nonZero := []float64{0, 0, 0.001}

	if !IsZeroVector(zero) {
		t.Error("Should detect zero vector")
	}
	if IsZeroVector(nonZero) {
		t.Error("Should not detect non-zero vector as zero")
	}
}
