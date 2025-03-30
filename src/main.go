package main

import (
	"rogue/datalayer"
	"rogue/domain"
	"rogue/presentation"
	"time"
)

func main() {
	// Создаём уровень и игрока
	session := domain.NewGameSession()

	// Создаём рендерер
	renderer := presentation.NewRenderer()

	// Основной игровой цикл
	for {
		// Рассчитываем туман войны
		fog := presentation.ComputeFogOfWar(session.Player, session.Levels[session.CurrentLevel], &session.Levels[session.CurrentLevel].Fog_corr)

		// Рендерим сцену
		renderer.Render(session, session.Levels[session.CurrentLevel], session.Player, fog)

		// Обрабатываем ввод
		presentation.HandleInput(renderer, session, session.Player, session.Levels[session.CurrentLevel])

		// Задержка для плавности
		time.Sleep(50 * time.Millisecond)

		if session.GameOver {
			renderer.GameOver()
			break
		}

		datalayer.Json_Save(session)
	}
}
