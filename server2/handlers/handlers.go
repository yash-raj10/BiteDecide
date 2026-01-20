package handlers

import (
	"net/http"
	"server2/engine"
	"server2/store"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	foodStore    *store.FoodStore
	sessionStore *store.SessionStore
	recommender  *engine.Recommender
}

// creates a new handler
func NewHandler(foodStore *store.FoodStore, sessionStore *store.SessionStore, recommender *engine.Recommender) *Handler {
	return &Handler{
		foodStore:    foodStore,
		sessionStore: sessionStore,
		recommender:  recommender,
	}
}

// handles /session
func (h *Handler) CreateSession(c *gin.Context) {
	sessionID := h.sessionStore.Create()
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
	})
}

// handles /recommendation
func (h *Handler) GetRecommendation(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	session := h.sessionStore.Get(sessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if session.IsCompleted() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session completed"})
		return
	}

	food := h.recommender.GetNextRecommendation(session)
	if food == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no more recommendations"})
		return
	}

	session.MarkSeen(food.ID)

	c.JSON(http.StatusOK, gin.H{
		"name":        food.Name,
		"description": food.Description,
	})
}

//request body for swipe
type SwipeRequest struct {
	SessionID string `json:"session_id"`
	FoodName  string `json:"food_name"`
	Action    string `json:"action"`
}

// handles /swipe
func (h *Handler) Swipe(c *gin.Context) {
	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Action != "left" && req.Action != "right" && req.Action != "super" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
		return
	}

	session := h.sessionStore.Get(req.SessionID)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if session.IsCompleted() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session completed"})
		return
	}

	food := h.foodStore.GetByName(req.FoodName)
	if food == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
		return
	}

	h.recommender.UpdateIntent(session, food, req.Action)

	if req.Action == "super" {
		session.Complete(req.FoodName)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
