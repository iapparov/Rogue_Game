package domain

type Character struct{
	Name string
	X, Y int
	MaxHealth int
	Health int
	Agility int
	Strength int
	Weapon *Item
	Backpack []*Item
}

// NewCharacter создаёт персонажа
func NewCharacter(name string, health, agility, strength, X, Y int) *Character {
	return &Character{
		X: X,
		Y: Y,
		Name:      name,
		Health:    health,
		MaxHealth: health,
		Agility: agility,
		Strength:  strength,
		Backpack:  []*Item{},
	}
}

// PickUpItem добавляет предмет в рюкзак
func (c *Character) PickUpItem(item *Item) {
	if len(c.Backpack) < 9 {
		c.Backpack = append(c.Backpack, item)
	}
}

func (c *Character) Move(X int, Y int, level *Level){
	c.X+=X
	c.Y+=Y
}