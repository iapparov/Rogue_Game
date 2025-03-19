package main

import (
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
		fog := presentation.ComputeFogOfWar(session.Player, session.Levels[session.CurrentLevel])

		// Рендерим сцену
		renderer.Render(session, session.Levels[session.CurrentLevel], session.Player, fog)

		// Обрабатываем ввод
		presentation.HandleInput(session.Player, session.Levels[session.CurrentLevel])

		// Задержка для плавности
		time.Sleep(50 * time.Millisecond)
	}
}