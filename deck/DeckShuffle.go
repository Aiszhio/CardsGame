package deck

import (
	"math/rand"
	"time"
)

func (deck *Deck) ShuffleAndTake() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len((*deck).Queue.Cards); i++ {
		temp := rand.Intn(len((*deck).Queue.Cards) - i)
		(*deck).Queue.Cards[i], (*deck).Queue.Cards[temp] = (*deck).Queue.Cards[temp], (*deck).Queue.Cards[i]
	}
	(*deck).Currents.Cards = (*deck).Queue.Cards[:8]
	(*deck).Queue.Cards = (*deck).Queue.Cards[8:]
}
