package deck

import (
	"fmt"
)

type DeckOutput interface {
	Print()
}

//func (d Deck) GetSuit(currSuit string) Deck {
//	var suitDeck Deck
//	for _, card := range d {
//		splitCard := strings.Split(card, " ")
//		for _, suit := range splitCard {
//			if suit == currSuit {
//				suitDeck = append(suitDeck, card)
//			}
//		}
//	}
//	return suitDeck
//}
//
//func (d Deck) PrintSuit() {
//	fmt.Printf("You've chosen a %s suit:\n%v\n", currSuit, d.GetSuit(currSuit))
//}

func (cur Currents) Print() {
	fmt.Println("Here's your deck: ")
	for _, card := range cur.Cards {
		fmt.Println(card.Rank, "of", card.Suit)
	}
}

func (que Queue) Print() {
	fmt.Println("Here's queue: ")
	for _, card := range que.Cards {
		fmt.Println(card.Rank, "of", card.Suit)
	}
}
