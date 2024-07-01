package inventory

type Inventory struct {
	Groups map[string]*Group
}

func NewInventory(name string) *Inventory {
	return &Inventory{
		Groups: make(map[string]*Group),
	}
}
