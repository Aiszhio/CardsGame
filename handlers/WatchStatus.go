package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClientData struct {
	UserHand []string `json:"hand"`
	AIHand   []string `json:"deckAI"`
	Table    []string `json:"table"`
}

func WatchStatus(c *gin.Context) {
	var data ClientData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prepareToFill := Status(data.UserHand, data.AIHand)

	var User, AI []string
	for idx, ok := range prepareToFill {
		if idx == 0 && !ok {
			User = FillHand(data.UserHand, &data.Table)
		} else if idx == 1 && !ok {
			AI = FillHand(data.AIHand, &data.Table)
		}
	}

	resp := gin.H{
		"handCards":  User,
		"aiCards":    AI,
		"tableCards": data.Table,
	}
	c.JSON(http.StatusOK, resp)

	sid := c.GetHeader("X-Session-ID")
	if len(data.UserHand) == 0 && len(data.Table) == 0 {
		push(sid, gin.H{"type": "gameOver", "winner": "player"})
	}
	if len(data.AIHand) == 0 && len(data.Table) == 0 {
		push(sid, gin.H{"type": "gameOver", "winner": "ai"})
	}
}

func Status(UserHand, AIHand []string) []bool {
	flags := make([]bool, 2)
	flags[0] = len(UserHand) == 8
	flags[1] = len(AIHand) == 8
	return flags
}

func FillHand(hand []string, table *[]string) []string {

	if len(hand) >= 8 || len(*table) == 0 {
		return hand
	}

	for len(hand) < 8 && len(*table) > 0 {
		hand = append(hand, (*table)[0])
		*table = (*table)[1:]
	}

	return hand
}
