package deck

import (
	"fmt"
	"strings"
)

func (d *Deck) GetSuit(currSuit string) Deck {
	var suitDeck Deck
	for _, card := range *d {
		splitCard := strings.Split(card, " ")
		for _, suit := range splitCard {
			if suit == currSuit {
				suitDeck = append(suitDeck, card)
			}
		}
	}
	return suitDeck
}

func (d *Deck) PrintSuit(currSuit string) {
	fmt.Printf("You've chosen a %s suit:\n%v\n", currSuit, d.GetSuit(currSuit))
}

func (d *Deck) Print() {
	fmt.Println("The deck consists of")
	for _, card := range *d {
		fmt.Printf("%v\n", card)
	}
}
