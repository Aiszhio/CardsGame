package handlers

import (
	"fmt"
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
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error from server": err.Error()})
	}

	fmt.Println("data", data)

	prepareToFill := Status(data.UserHand, data.AIHand)

	var User, AI []string

	for idx, value := range prepareToFill {
		if idx == 0 && value == false {
			User = FillHand(data.UserHand, &data.Table)
		} else if idx == 1 && value == false {
			AI = FillHand(data.AIHand, &data.Table)
		}
	}

	fmt.Println(User, AI)

	c.JSON(http.StatusOK, gin.H{"handCards": User, "aiCards": AI, "tableCards": data.Table})
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
