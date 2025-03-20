package domain

type ItemType string
const(
	Weapon ItemType = "Weapon"
	Potion ItemType = "Potion"
	Scroll ItemType = "Scroll"
	Treasure ItemType = "Treasure"
	Food ItemType = "Food"
)

type Item struct {
	Type       ItemType
	Subtype    string
	Health     int  // Для еды и эликсиров
	MaxHealth  int  // Для повышения максимального HP
	Agility  int  // Для свитков и эликсиров
	Strength   int  // Для оружия, свитков, эликсиров
	Cost       int  // Для сокровищ
}

// UseItem применяет предмет
func (c *Character) UseItem(item *Item) {
	switch item.Type {
	case Food:
		c.Health += item.Health
		if c.Health > c.MaxHealth {
			c.Health = c.MaxHealth
		}
	case Potion:
		c.Agility += item.Agility
	case Scroll:
		c.Strength += item.Strength
	case Weapon:
		c.Weapon = item
	}
}

