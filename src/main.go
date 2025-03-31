package main

import (
	"log"
	"rogue/datalayer"
	"rogue/domain"
	"rogue/presentation"
	"time"
)

func main() {

	// Создаём рендерер
	renderer := presentation.NewRenderer()

	// Создаём уровень и игрока
	session, score := menu(renderer)

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
			datalayer.Json_Record_Save(session, score)
			presentation.Records(renderer, score)
			renderer.GameOver()
			break
		}

		err:=datalayer.Json_Save(session)
		if err != nil {
			renderer.AddMessage(err.Error())
		}
	}
}

func menu(renderer *presentation.Renderer) (*domain.GameSession, *domain.LeaderBoards){
	var session *domain.GameSession
	var score domain.LeaderBoards
	var err error
	
	err = datalayer.Json_Record_Load(&score)
	if err != nil{
		log.Fatal(err.Error())
	}
	for i := 0;i <1;{
		menu, name:=presentation.Menu(renderer)
		switch menu {
		case 1:
			session = domain.NewGameSession(name)
			datalayer.Json_Free()
			i++
		case 2:
			session, err =datalayer.Json_Load()
			if err !=nil{
				log.Fatal(err.Error()) 
			}
			i++
		case 3:
			presentation.Records(renderer, &score)
		}
	}

	return session, &score
}