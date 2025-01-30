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

type Deck struct {
	Currents
	Queue
}

var Suits = []string{"Spades", "Hearts", "Diamonds", "Clubs"}
var Ranks = []string{"Ace", "King", "Queen", "Jack", "Ten", "Nine", "Eight", "Seven", "Six"}

func CreateDeck(c *gin.Context) {
	var deck Deck

	deck.Create()

	deck.Shuffle()

	var currentCards, queueCards []string

	for i := 0; i < len(deck.Currents.Cards); i++ {
		currentCards = append(currentCards, deck.Currents.Cards[i].Rank+deck.Currents.Cards[i].Suit)
	}

	for i := 0; i < len(deck.Queue.Cards); i++ {
		queueCards = append(queueCards, deck.Queue.Cards[i].Rank+deck.Queue.Cards[i].Suit)
	}

	c.JSON(200, gin.H{
		"deck":  deck.Currents.Cards,
		"queue": deck.Queue.Cards,
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
		temp := rand.Intn(len((*deck).Queue.Cards) - i)
		(*deck).Queue.Cards[i], (*deck).Queue.Cards[temp] = (*deck).Queue.Cards[temp], (*deck).Queue.Cards[i]
	}
	(*deck).Currents.Cards = (*deck).Queue.Cards[:8]
	(*deck).Queue.Cards = (*deck).Queue.Cards[8:]
}
