package handlers

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type Card struct {
	Rank string
	Suit string
}

type Currents struct {
	Cards []Card
}

type Queue struct {
	Cards []Card
}

type Table struct {
	Cards []Card
}

type Deck struct {
	Currents
	Queue
	Table
}

var Suits = []string{"Spades", "Hearts", "Diamonds", "Clubs"}
var Ranks = []string{"Ace", "King", "Queen", "Jack", "Ten", "Nine", "Eight", "Seven", "Six"}

func CreateDeck(c *gin.Context) {
	var deck Deck
	var deckAI Deck

	deck.Create()

	deck.Shuffle()

	CardsDistribution(&deck, &deckAI)

	c.JSON(200, gin.H{
		"deck":   deck.Currents.Cards,
		"deckAI": deckAI.Currents.Cards,
		"queue":  deck.Queue.Cards,
	})
}

func (deck *Deck) Create() {
	for _, suit := range Suits {
		for _, rank := range Ranks {
			(*deck).Queue.Cards = append((*deck).Queue.Cards, Card{rank, suit})
		}
	}
}

func (deck *Deck) Shuffle() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len((*deck).Queue.Cards); i++ {
		temp := rand.Intn(len((*deck).Queue.Cards))
		(*deck).Queue.Cards[i], (*deck).Queue.Cards[temp] = (*deck).Queue.Cards[temp], (*deck).Queue.Cards[i]
	}
}

func CardsDistribution(player, AI *Deck) {
	player.Currents.Cards = player.Queue.Cards[:8]
	player.Queue.Cards = player.Queue.Cards[8:]
	AI.Currents.Cards = player.Queue.Cards[:8]
	player.Queue.Cards = player.Queue.Cards[8:]
}
