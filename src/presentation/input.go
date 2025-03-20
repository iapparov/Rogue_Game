package presentation

import (
	"rogue/domain"

	"github.com/rthornton128/goncurses"
)

// HandleInput обрабатывает нажатия клавиш
func HandleInput(r *Renderer, session *domain.GameSession, player *domain.Character, level *domain.Level) {
		ch := r.window.GetChar()

		switch ch {
		case 'w', 'ц':
			player.Move(0, -1, level)
		case 's', 'ы':
			player.Move(0, 1, level)
		case 'a', 'ф':
			player.Move(-1, 0, level)
		case 'd', 'в':
			player.Move(1, 0, level)
		case 'q', 'й':
			session.EndGame()
		case goncurses.KEY_ENTER, '\n':
			if player.NextLevel(level){
				session.NextLevel()
			}
		}

}