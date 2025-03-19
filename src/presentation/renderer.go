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

	// Отрисовываем стены, пол и коридоры
	for _, room := range level.Rooms {
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				if fogOfWar[domain.Point{X: x, Y: y}] {
					r.window.MovePrint(y, x, string(UnknownChar))
				} else {
					r.window.MovePrint(y, x, string(FloorChar))
				}
			}
		}
	}

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