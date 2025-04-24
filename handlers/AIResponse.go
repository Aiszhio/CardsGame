package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserRequest struct {
	Selected []string `json:"selected"`
	HandAI   []string `json:"handAI"`
}

type ResponseToUser struct {
	SelectedAI []string `json:"selectedAI"`
	HandAI     []string `json:"handAI"`
}

func AIResponse(c *gin.Context) {
	var req UserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strongerCards := req.CompareCards()

	var resp ResponseToUser

	if len(strongerCards) == 0 || len(strongerCards) != len(req.Selected) {
		resp = ResponseToUser{
			SelectedAI: nil,
			HandAI:     append(req.HandAI, req.Selected...),
		}
	} else {
		resp = ResponseToUser{
			SelectedAI: strongerCards,
			HandAI:     RemoveCards(strongerCards, req.HandAI),
		}
	}

	if len(resp.HandAI) == 0 {
		push(c.GetHeader("X-Session-ID"), gin.H{
			"type":   "gameOver",
			"winner": "ai",
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (request *UserRequest) CompareCards() []string {
	var selectedCandidates []string

	remainingHand := make([]string, len(request.HandAI))
	copy(remainingHand, request.HandAI)

	for _, userCard := range request.Selected {
		userParts := strings.Split(userCard, " ")
		if len(userParts) < 3 {
			continue
		}
		userRank := FindRank(userParts[0])
		userSuit := userParts[2]

		var bestCandidate string
		var bestCandidateIndex int = -1
		var bestCandidateRank int

		for i, aiCard := range remainingHand {
			aiParts := strings.Split(aiCard, " ")
			if len(aiParts) < 3 {
				continue
			}
			aiRank := FindRank(aiParts[0])
			aiSuit := aiParts[2]

			if aiSuit == userSuit && aiRank > userRank {
				if bestCandidate == "" || aiRank < bestCandidateRank {
					bestCandidate = aiCard
					bestCandidateRank = aiRank
					bestCandidateIndex = i
				}
			}
		}

		if bestCandidate != "" {
			selectedCandidates = append(selectedCandidates, bestCandidate)
			remainingHand = append(remainingHand[:bestCandidateIndex], remainingHand[bestCandidateIndex+1:]...)
		} else {
			return nil
		}
	}

	if len(selectedCandidates) == len(request.Selected) {
		return selectedCandidates
	}
	return nil
}

func RemoveCards(cardsToRemove []string, hand []string) []string {
	remaining := make([]string, len(hand))
	copy(remaining, hand)
	for _, card := range cardsToRemove {
		for i, hcard := range remaining {
			if card == hcard {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return remaining
}

func FindRank(Rank string) int {
	switch Rank {
	case "Ace":
		return 10
	case "King":
		return 9
	case "Queen":
		return 8
	case "Jack":
		return 7
	case "Ten":
		return 6
	case "Nine":
		return 5
	case "Eight":
		return 4
	case "Seven":
		return 3
	case "Six":
		return 2
	default:
		return 0
	}
}
