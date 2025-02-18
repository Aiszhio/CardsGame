package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GameState struct {
	HandCards  []string `json:"handCards"`
	AICards    []string `json:"aiCards"`
	TableCards []string `json:"tableCards"`
}

func SaveGameState(c *gin.Context) {
	var state GameState
	if err := c.ShouldBindJSON(&state); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sb strings.Builder

	sb.WriteString("=== User Hand ===\n")
	for i, card := range state.HandCards {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, card))
	}

	sb.WriteString("\n=== AI Hand ===\n")
	for i, card := range state.AICards {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, card))
	}

	sb.WriteString("\n=== Table Cards ===\n")
	for i, card := range state.TableCards {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, card))
	}

	content := sb.String()

	filename := fmt.Sprintf("gamestate_%d.txt", time.Now().UnixNano())
	filePath := filepath.Join("static", "saved", filename)

	if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	downloadURL := fmt.Sprintf("/static/saved/%s", filename)
	c.JSON(http.StatusOK, gin.H{"downloadUrl": downloadURL})
}
