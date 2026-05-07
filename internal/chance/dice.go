package chance

import (
	"math/rand/v2"
)

type ThrowResult struct {
	Total int   `json:"total"`
	Rolls []int `json:"rolls"`
}

type Die struct {
	Sided int
}

// Roll returns the simulated result of rolling the die
func (d Die) Roll() int {
	return rand.IntN(d.Sided) + 1
}

type Cup struct {
	dice []Die
}

// NewCup creates a cupful of dice with the number of sides
// equal to sided, and the number of dice equal to amt.
func NewCup(dice ...Die) *Cup {
	d := make([]Die, 0, len(dice))
	d = append(d, dice...)
	return &Cup{dice: d}
}

// Shake is a semantic step that could be used to shuffle dice
// or reset state if needed in the future.
func (c *Cup) Shake() {
	rand.Shuffle(len(c.dice), func(i, j int) {
		c.dice[i], c.dice[j] = c.dice[j], c.dice[i]
	})
}

// Add adds one or more dice to the cup.
func (c *Cup) Add(dice ...Die) {
	c.dice = append(c.dice, dice...)
}

func (c *Cup) Throw() *ThrowResult {
	var total int
	rolls := make([]int, len(c.dice))
	for i, d := range c.dice {
		rolls[i] = d.Roll()
		total += rolls[i]
	}
	return &ThrowResult{Total: total, Rolls: rolls}
}
