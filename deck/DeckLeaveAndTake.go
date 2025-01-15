package deck

import "fmt"

func (deck *Deck) TakeCards() {

	count := 8 - len(deck.Currents.Cards)

	if len(deck.Currents.Cards) >= 8 || count == 0 {
		fmt.Println("Too many cards")
	} else if len(deck.Currents.Cards) == 0 {
		fmt.Println("You won!")
	}

	if len(deck.Queue.Cards) == 0 {
		fmt.Println("Queue is empty")
	}

	(*deck).Queue.Cards = append((*deck).Queue.Cards[count:])

}

func (deck *Deck) LeaveCards(cards ...Card) {

	if len(deck.Currents.Cards) == 0 {
		fmt.Println("No cards left")
	} else if len(deck.Currents.Cards) < len(cards) {
		fmt.Println("Not enough cards left")
	}

	for i := 0; i < len(deck.Currents.Cards); i++ {
		if deck.Currents.Cards[i] == cards[i] {
			(*deck).Currents.Cards = append(deck.Currents.Cards[:i], deck.Currents.Cards[i+1:]...)
		}
	}

}
