package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Console struct {
	open  bool
	input []rune
}

func (c *Console) Toggle() {
	c.open = !c.open
	if !c.open {
		c.input = c.input[:0]
	}
}

func (c *Console) IsOpen() bool {
	return c.open
}

func (c *Console) Input() string {
	return string(c.input)
}

func (c *Console) ReadInput() string {
	for ch := rl.GetCharPressed(); ch != 0; ch = rl.GetCharPressed() {
		if ch == '`' || ch == '~' {
			continue
		}
		c.input = append(c.input, ch)
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && len(c.input) > 0 {
		c.input = c.input[:len(c.input)-1]
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		cmd := string(c.input)
		c.input = c.input[:0]
		return cmd
	}

	return ""
}
