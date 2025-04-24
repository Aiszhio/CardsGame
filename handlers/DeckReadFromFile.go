package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewDeckFromFile(filename string) (handCards []Card, aiCards []Card, tableCards []Card, err error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	lines := strings.Split(string(bs), "\n")
	var nonEmpty []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}

	var currentSection string
	for _, line := range nonEmpty {
		switch line {
		case "Стол", "Выброшенные карты", "Ваша Рука", "Рука ИИ":
			currentSection = line
			continue
		}
		cardStr := line
		if dotIndex := strings.Index(cardStr, "."); dotIndex != -1 {
			cardStr = strings.TrimSpace(cardStr[dotIndex+1:])
		}
		parts := strings.Split(cardStr, " of ")
		if len(parts) != 2 {
			continue
		}
		card := Card{
			Rank: parts[0],
			Suit: parts[1],
		}
		switch currentSection {
		case "Стол":
			tableCards = append(tableCards, card)
		case "Ваша Рука":
			handCards = append(handCards, card)
		case "Рука ИИ":
			aiCards = append(aiCards, card)
		}
	}

	return handCards, aiCards, tableCards, nil
}

func LoadDeckFromFileHandler(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр filename обязателен"})
		return
	}

	filePath := filepath.Join("static", "saved", filename)

	handCards, aiCards, tableCards, err := NewDeckFromFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if handCards == nil {
		handCards = []Card{}
	}
	if aiCards == nil {
		aiCards = []Card{}
	}
	if tableCards == nil {
		tableCards = []Card{}
	}

	c.JSON(http.StatusOK, gin.H{
		"handCards":  handCards,
		"aiCards":    aiCards,
		"tableCards": tableCards,
	})
}
