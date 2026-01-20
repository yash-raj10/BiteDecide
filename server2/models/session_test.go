package models

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	session := NewSession("test-id", 10)

	if session.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", session.ID)
	}

	if len(session.IntentVector) != 10 {
		t.Errorf("Expected intent vector size 10, got %d", len(session.IntentVector))
	}

	if session.Completed {
		t.Error("New session should not be completed")
	}

	if len(session.SeenFoods) != 0 {
		t.Error("New session should have no seen foods")
	}
}

func TestSessionMarkSeen(t *testing.T) {
	session := NewSession("test", 3)

	session.MarkSeen("food-1")
	session.MarkSeen("food-2")

	if !session.HasSeen("food-1") {
		t.Error("Should have seen food-1")
	}
	if !session.HasSeen("food-2") {
		t.Error("Should have seen food-2")
	}
	if session.HasSeen("food-3") {
		t.Error("Should not have seen food-3")
	}
}

func TestSessionUpdateIntent(t *testing.T) {
	session := NewSession("test", 3)

	newIntent := []float64{0.5, 0.5, 0}
	session.UpdateIntent(newIntent)

	intent := session.GetIntent()
	for i := range newIntent {
		if intent[i] != newIntent[i] {
			t.Errorf("Intent[%d] = %v, want %v", i, intent[i], newIntent[i])
		}
	}
}

func TestSessionGetIntentReturnsClone(t *testing.T) {
	session := NewSession("test", 3)
	session.UpdateIntent([]float64{1, 2, 3})

	intent := session.GetIntent()
	intent[0] = 999

	original := session.GetIntent()
	if original[0] == 999 {
		t.Error("GetIntent should return a copy, not the original slice")
	}
}

func TestSessionComplete(t *testing.T) {
	session := NewSession("test", 3)

	if session.IsCompleted() {
		t.Error("New session should not be completed")
	}

	session.Complete("Pizza")

	if !session.IsCompleted() {
		t.Error("Session should be completed")
	}

	if session.FinalChoice != "Pizza" {
		t.Errorf("Final choice should be 'Pizza', got '%s'", session.FinalChoice)
	}
}
