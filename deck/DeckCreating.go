package deck

type Suit string

type Rank string

type Card struct {
	Rank Rank
	Suit Suit
}

type Currents struct {
	Cards []Card
}

type Queue struct {
	Cards []Card
}

type Deck struct {
	Cards []Card
	Currents
	Queue
}

var suits = []Suit{"Spades", "Hearts", "Diamonds", "Clubs"}
var ranks = []Rank{"Ace", "King", "Queen", "Jack", "Ten", "Nine", "Eight", "Seven", "Six"}

func (d *Deck) Create() {
	for _, suit := range suits {
		for _, rank := range ranks {
			(*d).Cards = append((*d).Cards, Card{rank, suit})
		}
	}
}
