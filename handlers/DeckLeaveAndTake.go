package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RequestData struct {
	Hand     []string `json:"hand"`
	Table    []string `json:"table"`
	Selected []string `json:"selected"`
}

type ResponseData struct {
	Deck        []string `json:"deck"`
	Queue       []string `json:"queue"`
	GamingTable []string `json:"table"`
}

func TakeCards(c *gin.Context) {
	var data RequestData

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(data.Hand) == 8 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "У вас слишком много карт"})
	} else if len(data.Table) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "На столе больше нет карт"})
	}

	resp := ResponseData{
		Deck:  data.Hand,
		Queue: data.Table,
	}

	numToTake := 8 - len(resp.Deck)
	for i := 0; i < numToTake && len(resp.Queue) > 0; i++ {
		card := resp.Queue[0]
		resp.Deck = append(resp.Deck, card)
		resp.Queue = resp.Queue[1:]
	}

	c.JSON(http.StatusOK, gin.H{
		"deck":  resp.Deck,
		"queue": resp.Queue,
	})
}
func LeaveCards(c *gin.Context) {
	var req RequestData

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Hand) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "В руке нет карт!"})
		return
	}

	resp := ResponseData{Deck: req.Hand, Queue: req.Table}

	var differentCards []string

	for _, selCard := range req.Selected {
		for _, card := range req.Hand {
			if selCard == card {
				differentCards = append(differentCards, card)
				break
			}
		}
	}

	for i := 0; i < len(differentCards)-1; i += 2 {
		if strings.TrimSpace(strings.Split(differentCards[i], " ")[0]) == strings.TrimSpace(strings.Split(differentCards[i+1], " ")[0]) {
			resp.GamingTable = append(resp.GamingTable, differentCards[i], differentCards[i+1])
			continue
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Карты должны быть одной величины!"})
			return
		}
	}

	remainingCards := []string{}
	for _, card := range req.Hand {
		found := false
		for _, selectedCard := range differentCards {
			if card == selectedCard {
				found = true
				break
			}
		}
		if !found {
			remainingCards = append(remainingCards, card)
		}
	}
	resp.Deck = remainingCards

	c.JSON(http.StatusOK, gin.H{
		"deck":        resp.Deck,
		"gamingTable": resp.GamingTable,
		"queue":       resp.Queue,
	})
}
