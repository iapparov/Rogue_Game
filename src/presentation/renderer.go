package presentation

import (
	"github.com/rthornton128/goncurses"
	"rogue/domain"
)

// Символы для рендеринга
const (
	WallChar      = '#'
	FloorChar     = '.'
	CorridorChar  = '+'
	PlayerChar    = '@'
	EnemyChar     = 'E'
	ItemChar      = '*'
	UnknownChar   = ' ' // Туман войны
	DoorChar = 'D'
)

// Renderer отвечает за отрисовку игрового мира
type Renderer struct {
	window *goncurses.Window
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
					if room == level.EndRoom{
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

	// // Отрисовываем предметы
	// for _, item := range level.Items {
	// 	if fogOfWar[domain.Point{X: item.X, Y: item.Y}] {
	// 		continue
	// 	}
	// 	r.window.MovePrint(item.Y, item.X, string(ItemChar))
	// }

	// Обновляем экран
	r.window.Refresh()
}

func (r *Renderer) GameOver(){
	goncurses.End()
}