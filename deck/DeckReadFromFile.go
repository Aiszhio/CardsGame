package deck

import (
	"fmt"
	"os"
	"strings"
)

func isCardValid(cardSlice []string) (error, []Card) {

	var cards []Card

	for i := 0; i < len(cardSlice)-2; i++ {
		if cardSlice[i] == "is" {
			cards = append(cards, Card{
				Rank: cardSlice[i+1],
				Suit: cardSlice[i+2],
			})
		}
	}

	return nil, cards
}

func NewDeckFromFile(filename string) Deck {
	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	strSlice := strings.Fields(string(bs))

	cards := getCardsFromFile(strSlice)

	var allCards []Card
	var currentDeck []Card

	for i := 0; i < len(cards); i++ {
		currentDeck = append(currentDeck, Card{
			Rank: cards[i].Rank,
			Suit: cards[i].Suit,
		})
	}

	for _, suit := range Suits {
		for _, rank := range Ranks {
			allCards = append(allCards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}

	for _, card := range cards {
		for idx, allCard := range allCards {
			if card.Suit == allCard.Suit && card.Rank == allCard.Rank {
				allCards = append(allCards[:idx], allCards[idx+1:]...)
			}
		}
	}

	deck := Deck{
		Currents{currentDeck},
		Queue{allCards},
	}

	fmt.Println("Here are your cards: \n", deck.Currents, "\n", "And the queue consists of \n", deck.Queue)

	return deck
}

func getCardsFromFile(cardSlice []string) []Card {

	var count int

	for idx, card := range cardSlice {
		if card == "1" {
			count = idx
			break
		}
	}

	cardSlice = cardSlice[count:]

	err, valuedCards := isCardValid(cardSlice)
	if err != nil {
		fmt.Println(err)
	}

	return valuedCards
}
