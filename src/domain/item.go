package domain

import(
	"math/rand"
)
/*  + сокровища (имеют стоимость, накапливаются и влияют на итоговый рейтинг, можно получить только при победе над монстром);
  + еда (восстанавливает здоровье на некоторую величину);
  + эликсиры (временно повышают одну из характеристик: ловкость, силу, максимальное здоровье);
  + свитки (постоянно повышают одну из характеристик: ловкость, силу, максимальное здоровье);
  + оружие (имеют характеристику силы, при использовании оружия меняется формула вычисления наносимого урона).
*/
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
	X,Y int
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
		c.Strength += item.Strength
		c.MaxHealth += item.MaxHealth
	case Scroll:
		c.Strength += item.Strength
		c.Agility += item.Agility
		c.MaxHealth += item.MaxHealth
	case Weapon:
		c.Weapon = item // допилить добавление у бавление статов
	}
}

func NewItem(name ItemType, subtype string, health, maxhealth, agility, strength, cost, x,y int) *Item{
	return &Item{
		Type: name, 
		Subtype: subtype,
		Health: health,
		MaxHealth: maxhealth,
		Agility: agility,
		Strength: strength,
		Cost: cost,
		X: x,
		Y: y,
	}
}

func GenerateItem(level *Level, current_level int) []*Item{
	


	gen_item := rand.Intn(7)+3

	items := make([]*Item, gen_item)

	itemsTypes := []*Item{
		NewItem(Weapon, "SuperSword", 0, 0, 0, 5+5*(current_level/5), 0, 0, 0),
		NewItem(Potion, "Strength Potion", 0, 0, 0, 5+5*(current_level/5), 0, 0, 0),
		NewItem(Potion, "Agility Potion", 0, 0, 5+5*(current_level/5), 0, 0, 0, 0),
		NewItem(Potion, "MaxHealth Potion", 5+5*(current_level/5), 5+5*(current_level/5), 0, 0, 0, 0, 0),
		NewItem(Scroll, "Strength Scroll", 0, 0, 0, 5+5*(current_level/5), 0, 0, 0),
		NewItem(Scroll, "Agility Scroll", 0, 0, 5+5*(current_level/5), 0, 0, 0, 0),
		NewItem(Scroll, "MaxHealth Scroll", 5+5*(current_level/5), 5+5*(current_level/5), 0, 0, 0, 0, 0),
		NewItem(Food, "Fish", 5+5*(current_level/5), 0, 0, 0, 0, 0, 0),
	}

	// Генерируем случайные предметы
	for i := 0; i < gen_item; i++ {
		randomIndex := rand.Intn(len(itemsTypes)) // Выбираем случайный предмет из списка
		item := itemsTypes[randomIndex]
		items[i] = item
		items[i].X = level.Rooms[i].X+2
		items[i].Y = level.Rooms[i].Y+2
	}

	
	return items
}