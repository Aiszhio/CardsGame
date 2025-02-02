package handlers

import (
	"fmt"
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
	Deck  []string `json:"deck"`
	Queue []string `json:"queue"`
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

	fmt.Println(len(resp.Deck), len(resp.Queue), 8-len(resp.Deck))

	for i := 0; i < 8-len(resp.Deck); i++ {
		resp.Deck = append(resp.Deck, resp.Queue[i])
		resp.Queue = append(resp.Queue[:i], resp.Queue[i+1:]...)
		fmt.Println(resp.Deck, resp.Queue)
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

	cards := make([]string, 0, len(req.Selected))

	resp := ResponseData{Deck: req.Hand, Queue: req.Table}

	for _, selCard := range req.Selected {
		for i, card := range req.Hand {
			if selCard == card {
				resp.Deck = append(resp.Deck[:i], resp.Deck[i+1:]...)
				break
			}
		}
	}

	for _, card := range req.Selected {
		trimmed := strings.TrimSpace(strings.Split(card, " ")[0])
		cards = append(cards, trimmed)
	}

	for i := 0; i < len(cards)-1; i++ {
		if cards[i] == cards[i+1] {
			continue
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Карты должны быть одной величины!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"deck":  resp.Deck,
		"queue": resp.Queue,
	})
}
