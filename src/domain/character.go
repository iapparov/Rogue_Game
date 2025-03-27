package domain

import "strconv"

type Character struct {
	Name          string
	X, Y          int
	MaxHealth     int
	Health        int
	Agility       int
	Strength      int
	Weapon        *Item
	Weapon_hand   bool
	Backpack      []*Item
	TreasureCount int
}

// NewCharacter создаёт персонажа
func NewCharacter(name string, health, agility, strength, X, Y int) *Character {
	return &Character{
		X:           X,
		Y:           Y,
		Name:        name,
		Health:      health,
		MaxHealth:   health,
		Agility:     agility,
		Strength:    strength,
		Backpack:    []*Item{},
		Weapon_hand: false,
	}
}

// PickUpItem добавляет предмет в рюкзак
func (c *Character) PickUpItem(level *Level) (int, ItemType) {
	for i := 0; i < len(level.Items); i++ {
		if level.Items[i].X == c.X && level.Items[i].Y == c.Y {
			if level.Items[i].Type == Treasure {
				// Для сокровищ не требуется места в рюкзаке
				c.TreasureCount += level.Items[i].Cost
				level.Items[i].X = -1
				level.Items[i].Y = -1
				return 2, level.Items[i].Type
			}
			if len(c.Backpack) < 9 {
				c.Backpack = append(c.Backpack, level.Items[i])
				level.Items[i].X = -1
				level.Items[i].Y = -1
				return 1, level.Items[i].Type
			} else {
				return 0, level.Items[i].Type
			}
		}
	}
	return -1, ""
}

func (c *Character) UseH(firstch, ch string, level *Level) int {
	numb, err := strconv.Atoi(ch)

	// Проверяем, существует ли элемент в рюкзаке
	if numb <= 0 || numb > len(c.Backpack) || err != nil {
		return -2 // Неверный индекс
	} else {
		switch firstch {
		case "h":
			if c.Backpack[numb-1].Type != Weapon {
				return -2
			} else {
				if c.Weapon_hand {
					level.Items = append(level.Items, c.Backpack[numb-1])
					level.Items[len(level.Items)-1].X = c.X
					level.Items[len(level.Items)-1].Y = c.Y + 1
				}
				c.UseItem(c.Backpack[numb-1])
				c.Backpack = append(c.Backpack[:numb-1], c.Backpack[numb:]...)
				return 200
			}
		case "j":
			if c.Backpack[numb-1].Type != Food {
				return -2
			} else {
				c.UseItem(c.Backpack[numb-1])
				c.Backpack = append(c.Backpack[:numb-1], c.Backpack[numb:]...)
				return 200
			}
		case "k":
			if c.Backpack[numb-1].Type != Potion {
				return -2
			} else {
				c.UseItem(c.Backpack[numb-1])
				c.Backpack = append(c.Backpack[:numb-1], c.Backpack[numb:]...)
				return 200
			}
		case "e":
			if c.Backpack[numb-1].Type != Scroll {
				return -2
			} else {
				c.UseItem(c.Backpack[numb-1])
				c.Backpack = append(c.Backpack[:numb-1], c.Backpack[numb:]...)
				return 200
			}
		}
	}
	return -2
}

func (c *Character) NextLevel(level *Level) bool {
	if c.X == level.EndRoom.DoorX && c.Y == level.EndRoom.DoorY {
		return true
	}
	return false
}

func (c *Character) Move(dx, dy int, level *Level) {
	newX, newY := c.X+dx, c.Y+dy

	// Проверяем, находится ли новая позиция в комнате или коридоре
	if isWalkable(newX, newY, level) && !isEnemyAt(newX, newY, level) {
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

func isEnemyAt(x, y int, level *Level) bool {
	for _, enemy := range level.Enemies {
		if enemy.X == x && enemy.Y == y {
			return true
		}
	}
	return false
}
