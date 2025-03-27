package presentation

import (
	"rogue/domain"

	"github.com/rthornton128/goncurses"
)

// Символы для рендеринга
const (
	WallChar     = '|'
	FloorChar    = '.'
	CorridorChar = '+'
	PlayerChar   = '@'
	EnemyChar    = 'E'
	ItemChar     = '*'
	UnknownChar  = ' ' // Туман войны
	DoorChar     = 'D'
)

// Renderer отвечает за отрисовку игрового мира
type Renderer struct {
	window   *goncurses.Window
	messages []string // Буфер сообщений
	backpack bool
}

// NewRenderer создаёт новый рендерер
func NewRenderer() *Renderer {
	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	stdscr.Keypad(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)

	return &Renderer{window: stdscr}
}

// Render отрисовывает уровень
func (r *Renderer) Render(session *domain.GameSession, level *domain.Level, player *domain.Character, fogOfWar map[domain.Point]bool) {
	r.window.Clear()

	// Отрисовываем коридоры
	for _, corridor := range level.Corridors {
		for _, point := range corridor.Path {
			if fogOfWar[point] {
				r.window.MovePrint(point.Y, point.X, string(UnknownChar))
			} else {
				r.window.MovePrint(point.Y, point.X, string(CorridorChar))
			}
		}
	}

	// Отрисовываем стены
	for _, room := range level.Rooms {
		for _, wall := range room.Walls {
			if fogOfWar[wall] {
				r.window.MovePrint(wall.Y, wall.X, string(UnknownChar))
			} else {
				r.window.MovePrint(wall.Y, wall.X, string(WallChar))
			}
		}
	}

	// Отрисовываем пол
	for _, room := range level.Rooms {
		for y := room.Y + 1; y < room.Y+room.Height-1; y++ {
			for x := room.X + 1; x < room.X+room.Width-1; x++ {
				if fogOfWar[domain.Point{X: x, Y: y}] {
					r.window.MovePrint(y, x, string(UnknownChar))
				} else {
					r.window.MovePrint(y, x, string(FloorChar))
					if room == level.EndRoom {
						r.window.MovePrint(level.EndRoom.DoorY, level.EndRoom.DoorX, string(DoorChar))
					}
				}
			}
		}
	}

	// Отрисовываем персонажа
	r.window.MovePrint(player.Y, player.X, string(PlayerChar))

	// // Отрисовываем противников
	// for _, enemy := range level.Enemies {
	// 	if fogOfWar[domain.Point{X: enemy.X, Y: enemy.Y}] {
	// 		continue
	// 	}
	// 	r.window.MovePrint(enemy.Y, enemy.X, string(EnemyChar))
	// }

	// Отрисовываем предметы
	for _, item := range level.Items {
		if fogOfWar[domain.Point{X: item.X, Y: item.Y}] {
			continue
		}
		r.window.MovePrint(item.Y, item.X, string(ItemChar))
	}

	// Прописываем статы

	r.window.MovePrint(30, 0, "MaxHealth: ", player.MaxHealth)
	r.window.MovePrint(31, 0, "Health: ", player.Health)
	r.window.MovePrint(32, 0, "Agility: ", player.Agility)
	r.window.MovePrint(33, 0, "Strength: ", player.Strength)
	r.window.MovePrint(34, 0, "Curren Level: ", session.CurrentLevel+1)

	// Вывод последних сообщений
	startY := 35
	for i, msg := range r.messages {
		r.window.MovePrint(startY+i, 0, msg)
	}

	if r.backpack {
		r.BackPack(player)
	}

	// Обновляем экран
	r.window.Refresh()
}

func (r *Renderer) GameOver() {
	goncurses.End()
}

// AddMessage добавляет новое сообщение в буфер
func (r *Renderer) AddMessage(msg string) {
	// Ограничиваем количество сообщений (например, 3)
	if len(r.messages) >= 3 {
		r.messages = r.messages[1:] // Удаляем самое старое сообщение
	}
	r.messages = append(r.messages, msg)
}

func (r *Renderer) TakeSomething(flag int, message domain.ItemType) {
	var msg string
	switch flag {
	case 1:
		msg = string(message) + " now in backpack"
		r.AddMessage(msg)
	case 0:
		msg = "Backpack is full. Can't take " + string(message)
		r.AddMessage(msg)
	case -2:
		r.AddMessage("Can't use it")
	}
}

func (r *Renderer) BackPack(player *domain.Character) {
	backpack := player.Backpack
	for i, item := range backpack {
		r.window.MovePrint(1+i, 50, (i + 1), ") ", item.Subtype, "(", "Ag +", item.Agility, " He+", item.Health, " MaxHe+", item.MaxHealth, " Str+", item.Strength, ")")
	}
	if player.Weapon_hand {
		r.window.MovePrint(0, 50, 0, ") ", player.Weapon.Subtype, "(Str +", player.Weapon.Strength, ")")
	}
}
