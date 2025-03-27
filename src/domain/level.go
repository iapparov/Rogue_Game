package domain

import (
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
	Rooms     []*Room
	Items []*Item
	// Enemies		[]*Enemy
	Corridors []*Corridor
	StartRoom *Room
	EndRoom   *Room
	Fog_corr  map[Point]bool
}

// Room представляет комнату
type Room struct {
	X, Y      int // Верхний левый угол
	Width     int
	Height    int
	Walls     []Point   // Координаты стен комнаты
	Connected []*Room   // Связанные комнаты
	DoorX, DoorY int
}

// Corridor соединяет две комнаты
type Corridor struct {
	Path []Point
}

// Point - координаты клетки
type Point struct {
	X, Y int
}

// GenerateLevel создаёт случайный уровень
func GenerateLevel(depth int) *Level {
	level := &Level{}
	level.Fog_corr = make(map[Point]bool)
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
	level.EndRoom.DoorX = level.EndRoom.X+1
	level.EndRoom.DoorY = level.EndRoom.Y+1

	// 5. Генерируем предметы в каждой комнате 
	level.Items = GenerateItem(level, depth)
	

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
		walls = append(walls, Point{X: x, Y: room.Y})                     // Верхняя стена
		walls = append(walls, Point{X: x, Y: room.Y + room.Height - 1}) // Нижняя стена
	}

	// Левая и правая стенки
	for y := room.Y; y < room.Y+room.Height; y++ {
		walls = append(walls, Point{X: room.X, Y: y})                     // Левая стена
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

// Вспомогательная функция для вычисления абсолютного значения
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}