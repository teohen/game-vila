package inventory

const (
	Wood  = "wood"
	Food  = "food"
	Hide  = "hide"
	Stone = "stone"
)

type Inventory struct {
	wood  int
	meat  int
	hide  int
	stone int
	iron  int
}

var inventory Inventory

func NewInventory() *Inventory {
	inventory = Inventory{
		wood:  0,
		meat:  0,
		hide:  0,
		iron:  0,
		stone: 0,
	}
	return &inventory
}

func (i *Inventory) AddToInventory(res string, amount int) {
	switch res {
	case Wood:
		i.wood += amount
	case Food:
		i.meat += amount
	case Hide:
		i.hide += amount
	case Stone:
		i.stone += amount
	}
}

func Get() *Inventory {
	return &inventory
}
