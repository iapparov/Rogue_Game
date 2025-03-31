package domain

// GameSession управляет состоянием игры
type GameSession struct {
	CurrentLevel  int        `json:"CurrentLevel"`
	Player        *Character `json:"Player"`
	Levels        []*Level   `json:"Levels"`
	GameOver      bool       `json:"GameOver"`
	TreasureCount int        `json:"TreasureCount"`
}

type LeaderBoards struct{
	Record []*Records `json:"Records"`
}

type Records struct{
	Name string `json:"Name"`
	Record int `json:"Record"`
}

// NewGameSession создаёт новую игру
func NewGameSession(name string) *GameSession {
	player := NewCharacter(name, 100, 10, 5, 0, 0)
	session := &GameSession{
		CurrentLevel: 0,
		Player:       player,
		Levels:       make([]*Level, 21),
	}

	// Генерируем 21 уровень
	for i := 0; i < 21; i++ {
		session.Levels[i] = GenerateLevel(i + 1)
	}

	session.Player.X = session.Levels[0].StartRoom.X + 1
	session.Player.Y = session.Levels[0].StartRoom.Y + 1

	return session
}

// NextLevel переносит игрока на следующий уровень
func (g *GameSession) NextLevel() {
	if g.CurrentLevel < 21 {
		g.CurrentLevel++
	}
	g.Player.X = g.Levels[g.CurrentLevel].StartRoom.X + 1
	g.Player.Y = g.Levels[g.CurrentLevel].StartRoom.Y + 1
}

// EndGame завершает игру
func (g *GameSession) EndGame() {
	g.GameOver = true
}

