package item

const WOOD_WEIGHT = 5

type Wood struct {
	ID     string
	weight int
}

func NewWood() *Wood {
	return &Wood{
		ID:     "wood",
		weight: WOOD_WEIGHT,
	}
}
