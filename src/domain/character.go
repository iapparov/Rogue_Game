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

func (c *Character) NextLevel(level *Level) bool{
	if c.X == level.EndRoom.DoorX && c.Y == level.EndRoom.DoorY{
		return true
	}
	return false
}

func (c *Character) Move(dx, dy int, level *Level) {
	newX, newY := c.X+dx, c.Y+dy

	// Проверяем, находится ли новая позиция в комнате или коридоре
	if isWalkable(newX, newY, level) {
		c.X = newX
		c.Y = newY
	}
}

// Функция проверяет, можно ли пройти в указанную клетку
func isWalkable(x, y int, level *Level) bool {
	// Проверяем, находится ли в комнате
	for _, room := range level.Rooms {
		if x > room.X && x < room.X+room.Width-1 &&
			y > room.Y && y < room.Y+room.Height-1 {
			return true
		}
	}

	// Проверяем, находится ли в коридоре
	for _, corridor := range level.Corridors {
		for _, point := range corridor.Path {
			if point.X == x && point.Y == y {
				return true
			}
		}
	}

	return false
}