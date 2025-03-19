package presentation

import (
	"github.com/rthornton128/goncurses"
	"rogue/domain"
)

// HandleInput обрабатывает нажатия клавиш
func HandleInput(player *domain.Character, level *domain.Level) {
	for {
		ch := goncurses.StdScr().GetChar()

		switch ch {
		case 'w':
			player.Move(0, -1, level)
		case 's':
			player.Move(0, 1, level)
		case 'a':
			player.Move(-1, 0, level)
		case 'd':
			player.Move(1, 0, level)
		case 'q':
			return // Выход из игры
		}
	}
}