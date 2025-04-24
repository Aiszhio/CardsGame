package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResumePayload struct {
	HandCards  []string `json:"handCards"`
	AICards    []string `json:"aiCards"`
	TableCards []string `json:"tableCards"`
}

var stateBySession = make(map[string]ResumePayload)

func ResumeGame(c *gin.Context) {
	var p ResumePayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sid := c.GetHeader("X-Session-ID")
	stateBySession[sid] = p

	push(sid, gin.H{
		"type":   "state",
		"deck":   p.HandCards,
		"deckAI": p.AICards,
		"queue":  p.TableCards,
	})
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
