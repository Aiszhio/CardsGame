package deck

import (
	"math/rand"
	"time"
)

func (d *Deck) ShuffleAndTake() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len((*d).Cards); i++ {
		temp := rand.Intn(len((*d).Cards) - i)
		(*d).Cards[i], (*d).Cards[temp] = (*d).Cards[temp], (*d).Cards[i]
	}
	(*d).Currents.Cards = (*d).Cards[:8]
	(*d).Queue.Cards = (*d).Cards[8:]
}
