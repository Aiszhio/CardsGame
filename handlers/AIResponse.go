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
	}

	strongerCards := req.CompareCards()

	var resp ResponseToUser

	if len(strongerCards) == 0 {
		resp = ResponseToUser{
			SelectedAI: nil,
			HandAI:     append(req.HandAI, req.Selected...),
		}
	} else {
		resp = ResponseToUser{
			SelectedAI: strongerCards,
			HandAI:     ThrowCards(strongerCards, req.HandAI),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (request *UserRequest) CompareCards() []string {

	var strongerCards []string

	for i := 0; i < len(request.Selected); i++ {
		for j := 0; j < len(request.HandAI); j++ {
			result := func() bool {
				if strings.Split(request.Selected[i], " ")[2] == strings.Split(request.HandAI[j], " ")[2] {
					return true
				} else {
					return false
				}
			}

			if result() == true && FindRank(strings.Split(request.Selected[i], " ")[0]) < FindRank(strings.Split(request.HandAI[j], " ")[0]) {
				strongerCards = append(strongerCards, request.HandAI[j])
			}
		}
	}

	if len(strongerCards) == 0 {
		return nil
	}

	return strongerCards
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
func ThrowCards(cardsToThrow, hand []string) []string {
	if len(cardsToThrow) == 0 {
		return hand
	}

	weakestCard := cardsToThrow[0]
	weakestRank := FindRank(strings.Split(weakestCard, " ")[0])
	for _, card := range cardsToThrow {
		rank := FindRank(strings.Split(card, " ")[0])
		if rank < weakestRank {
			weakestRank = rank
			weakestCard = card
		}
	}

	for i, card := range hand {
		if card == weakestCard {
			newHand := append([]string{}, hand[:i]...)
			newHand = append(newHand, hand[i+1:]...)
			return newHand
		}
	}

	return hand
}
