package deck

import (
	"fmt"
)

type Output interface {
	Print()
}

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
