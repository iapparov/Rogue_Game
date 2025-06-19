package domain

import (
	"fmt"
	"math/rand"
)

// Размер карты
const (
	LevelWidth  = 50
	LevelHeight = 25
	MinRoomSize = 4
	MaxRoomSize = 8
)

// Level представляет уровень игры
type Level struct {
	Rooms     []*Room        `json:"Rooms"`
	Items     []*Item        `json:"Items"`
	Enemies   []*Enemy       `json:"Enemies"`
	Corridors []*Corridor    `json:"Corridors"`
	StartRoom *Room          `json:"StartRoom"`
	EndRoom   *Room          `json:"EndRoom"`
	Fog_corr  map[string]bool `json:"Fog_corr"`
}

// Room представляет комнату
type Room struct {
	X         int     `json:"X"`
	Y         int     `json:"Y"` // Верхний левый угол
	Width     int     `json:"Width"`
	Height    int     `json:"Height"`
	Walls     []Point `json:"Walls"`     // Координаты стен комнаты
	Connected []*Room `json:"Connected"` // Связанные комнаты
	DoorX     int     `json:"DoorX"`
	DoorY     int     `json:"DoorY"`
}

// Corridor соединяет две комнаты
type Corridor struct {
	Path []Point `json:"Path"`
}

// Point - координаты клетки
type Point struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

// GenerateLevel создаёт случайный уровень
func GenerateLevel(depth int) *Level {
	level := &Level{}
	level.Fog_corr = make(map[string]bool)
	sections := divideIntoSections(LevelWidth, LevelHeight)

	// 1. Генерируем комнаты в секциях
	for _, section := range sections {
		room := generateRoom(section)
		room.Walls = generateWalls(room) // Генерация стен
		level.Rooms = append(level.Rooms, room)
	}

	// 2. Связываем комнаты (MST)
	level.Corridors = connectRooms(level.Rooms)

	// 3. Назначаем стартовую и конечную комнаты
	level.StartRoom = level.Rooms[rand.Intn(len(level.Rooms))]
	level.EndRoom = level.Rooms[rand.Intn(len(level.Rooms))]
	for level.StartRoom == level.EndRoom { // Гарантируем разные комнаты
		level.EndRoom = level.Rooms[rand.Intn(len(level.Rooms))]
	}

	// 4. Создаем дверь для перехода на следующий уровень
	level.EndRoom.DoorX = level.EndRoom.X + 1
	level.EndRoom.DoorY = level.EndRoom.Y + 1

	// 5. Генерируем предметы в каждой комнате
	level.Items = GenerateItem(level, depth)
	level.SpawnEnemies()

	// 6. Разрываем циклические ссылки
	visited:=make(map[*Room]bool)
	RemoveCycles(level.Rooms, visited)
	
	return level

}

// Разбивает уровень на 9 секций
func divideIntoSections(width, height int) []Room {
	sections := []Room{}
	secWidth := width / 3
	secHeight := height / 3

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			sections = append(sections, Room{
				X:      col * secWidth,
				Y:      row * secHeight,
				Width:  secWidth,
				Height: secHeight,
			})
		}
	}
	return sections
}

// Генерирует случайную комнату в заданной секции
func generateRoom(section Room) *Room {
	w := rand.Intn(MaxRoomSize-MinRoomSize) + MinRoomSize
	h := rand.Intn(MaxRoomSize-MinRoomSize) + MinRoomSize
	x := section.X + rand.Intn(section.Width-w)
	y := section.Y + rand.Intn(section.Height-h)

	return &Room{X: x, Y: y, Width: w, Height: h}
}

// Генерирует стены вокруг комнаты
func generateWalls(room *Room) []Point {
	walls := []Point{}

	// Верхняя и нижняя стенки
	for x := room.X; x < room.X+room.Width; x++ {
		walls = append(walls, Point{X: x, Y: room.Y})                   // Верхняя стена
		walls = append(walls, Point{X: x, Y: room.Y + room.Height - 1}) // Нижняя стена
	}

	// Левая и правая стенки
	for y := room.Y; y < room.Y+room.Height; y++ {
		walls = append(walls, Point{X: room.X, Y: y})                  // Левая стена
		walls = append(walls, Point{X: room.X + room.Width - 1, Y: y}) // Правая стена
	}

	return walls
}

// Соединяет комнаты коридорами (алгоритм MST - Минимальное остовное дерево)
func connectRooms(rooms []*Room) []*Corridor {
	corridors := []*Corridor{}
	unvisited := make(map[*Room]bool)
	for _, r := range rooms {
		unvisited[r] = true
	}

	// Начинаем с первой комнаты
	current := rooms[0]
	unvisited[current] = false

	for len(unvisited) > 0 {
		// Находим ближайшую несоединённую комнату
		closestRoom, corridor := findClosestRoom(current, unvisited)
		if closestRoom == nil {
			break
		}

		// Соединяем их
		corridors = append(corridors, corridor)
		current.Connected = append(current.Connected, closestRoom)
		closestRoom.Connected = append(closestRoom.Connected, current)

		// Помечаем посещённой
		delete(unvisited, closestRoom)
		current = closestRoom
	}

	return corridors
}

// Поиск ближайшей комнаты и построение коридора
func findClosestRoom(current *Room, unvisited map[*Room]bool) (*Room, *Corridor) {
	var closestRoom *Room
	minDist := 99999
	var bestCorridor *Corridor

	for room := range unvisited {
		distance := abs(current.X-room.X) + abs(current.Y-room.Y)
		if distance < minDist {
			minDist = distance
			closestRoom = room
			bestCorridor = buildCorridor(current, closestRoom)
		}
	}
	return closestRoom, bestCorridor
}

// Создаёт коридор между двумя комнатами
func buildCorridor(r1, r2 *Room) *Corridor {
	path := []Point{}

	// Горизонтальное перемещение
	x1, y1 := r1.X+r1.Width/2, r1.Y+r1.Height/2
	x2, y2 := r2.X+r2.Width/2, r2.Y+r2.Height/2

	if x1 < x2 {
		for x := x1; x <= x2; x++ {
			path = append(path, Point{X: x, Y: y1})
		}
	} else {
		for x := x1; x >= x2; x-- {
			path = append(path, Point{X: x, Y: y1})
		}
	}

	// Вертикальное перемещение
	if y1 < y2 {
		for y := y1; y <= y2; y++ {
			path = append(path, Point{X: x2, Y: y})
		}
	} else {
		for y := y1; y >= y2; y-- {
			path = append(path, Point{X: x2, Y: y})
		}
	}

	return &Corridor{Path: path}
}

func (l *Level) SpawnEnemies() {
	for i := 0; i < len(l.Rooms); i++ {
		room := l.Rooms[i]
		x := room.X + 1 + rand.Intn(room.Width-2)
		y := room.Y + 1 + rand.Intn(room.Height-2)
		
		enemyType := getRandomEnemyType()
		enemy := NewEnemy(enemyType)
		enemy.X, enemy.Y = x, y

		l.Enemies = append(l.Enemies, enemy)
	}
}

func getRandomEnemyType() EnemyType {
	types := []EnemyType{Zombie, Vampire, Ghost, Ogre, SnakeMage}
	return types[rand.Intn(len(types))]
}
func (l *Level) GetEnemyAt(x, y int) *Enemy {
	for _, enemy := range l.Enemies {
		if abs(enemy.X-x) <= 1 && abs(enemy.Y-y) <= 1 {
			return enemy
		}
	}
	return nil
}

// RemoveEnemy удаляет врага из списка врагов уровня
func (l *Level) RemoveEnemy(enemy *Enemy) {
	for i, e := range l.Enemies {
		if e == enemy {
			treasure := enemy.DropTreasure()
			l.Items = append(l.Items, treasure)

			l.Enemies = append(l.Enemies[:i], l.Enemies[i+1:]...)
			break
		}
	}
}

// Вспомогательная функция для вычисления абсолютного значения
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// EnemiesAttack заставляет всех врагов атаковать игрока, если они рядом
func EnemiesAttack(player *Character, enemies []*Enemy, messages *[]string) {
	for _, enemy := range enemies {
		if isAdjacent(enemy.X, enemy.Y, player.X, player.Y) {
			if enemy.Attack(player) {
				damage := calculateDamage(&enemy.Character)
				*messages = append(*messages, fmt.Sprintf("%s attacks! Damage: %d", enemy.Name, damage))
			} else {
				*messages = append(*messages, fmt.Sprintf("%s missed!", enemy.Name))
			}
		}
	}
}

// isAdjacent проверяет, находится ли игрок рядом с врагом (по клеткам)
func isAdjacent(x1, y1, x2, y2 int) bool {
	return abs(x1-x2) <= 1 && abs(y1-y2) <= 1
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// FromString — создаёт Point из строки (например, "10,5")
func FromString(s string) (Point, error) {
	var p Point
	_, err := fmt.Sscanf(s, "%d,%d", &p.X, &p.Y)
	return p, err
}

func RemoveCycles(rooms []*Room, visited map[*Room]bool) {
	for _, room := range rooms {
		if visited[room] {
			continue
		}
		visited[room] = true

		for i := range room.Connected {
			if visited[room.Connected[i]] {
				room.Connected[i] = nil // Разрываем цикл
			}
		}
	}
}