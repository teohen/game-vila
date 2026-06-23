package inventory

const (
	Wood  = "wood"
	Food  = "food"
	Hide  = "hide"
	Stone = "stone"
)

type Inventory struct {
	Wood  int
	Meat  int
	Hide  int
	Stone int
	Iron  int
}

var inventory Inventory

func NewInventory() *Inventory {
	inventory = Inventory{
		Wood:  0,
		Meat:  0,
		Hide:  0,
		Iron:  0,
		Stone: 0,
	}
	return &inventory
}

func (i *Inventory) AddToInventory(res string, amount int) {
	switch res {
	case Wood:
		i.Wood += amount
	case Food:
		i.Meat += amount
	case Hide:
		i.Hide += amount
	case Stone:
		i.Stone += amount
	}
}

func Get() *Inventory {
	return &inventory
}
