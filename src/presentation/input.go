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
		flg, msg := player.PickUpItem(level)
		r.TakeSomething(flg, msg)
		r.backpack = false
	case 's', 'ы':
		player.Move(0, 1, level)
		flg, msg := player.PickUpItem(level)
		r.TakeSomething(flg, msg)
		r.backpack = false
	case 'a', 'ф':
		player.Move(-1, 0, level)
		flg, msg := player.PickUpItem(level)
		r.TakeSomething(flg, msg)
		r.backpack = false
	case 'd', 'в':
		player.Move(1, 0, level)
		flg, msg := player.PickUpItem(level)
		r.TakeSomething(flg, msg)
		r.backpack = false
	case 'q', 'й':
		session.EndGame()
	case 'h', 'j', 'k', 'e': //оружие
		r.backpack = true
		r.BackPack(player)
		flg := player.UseH(string(rune(ch)), string(rune(r.window.GetChar())), level) // Передаем символ в функцию
		r.TakeSomething(flg, "")
		r.backpack = false
	case 'f':
		enemy := level.GetEnemyAt(player.X, player.Y) // Проверяем, есть ли враг рядом
		if enemy != nil {
			success := player.Attack(&enemy.Character)
			if success {
				r.AddMessage("You hit " + enemy.Name + "!")
			} else {
				r.AddMessage("You missed!")
			}
			if enemy.Health <= 0 {
				r.AddMessage(enemy.Name + " is defeated!")
				level.RemoveEnemy(enemy)
			}
		} else {
			r.AddMessage("No enemy to attack!")
		}
		// Вызов атаки врагов после хода игрока
		domain.EnemiesAttack(player, level.Enemies, &r.messages)
		if player.Health <= 0 {
			session.EndGame()
		}

	case goncurses.KEY_ENTER, '\n':
		if player.NextLevel(level) && session.CurrentLevel == 21{
			session.EndGame()
		} else if player.NextLevel(level) {
			session.NextLevel()
		}
	}

}
