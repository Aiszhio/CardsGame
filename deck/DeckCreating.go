package deck

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

func (deck *Deck) Create() {
	for _, suit := range Suits {
		for _, rank := range Ranks {
			(*deck).Queue.Cards = append((*deck).Queue.Cards, Card{rank, suit})
		}
	}
}
