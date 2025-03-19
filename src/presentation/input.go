package presentation

import (
	"rogue/domain"
)

// HandleInput обрабатывает нажатия клавиш
func HandleInput(r *Renderer, session *domain.GameSession, player *domain.Character, level *domain.Level) {
		ch := r.window.GetChar()

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
			session.EndGame()
		}
}