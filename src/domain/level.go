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
	Corridors []*Corridor
	StartRoom *Room
	EndRoom   *Room
}

// Room представляет комнату
type Room struct {
	X, Y      int // Верхний левый угол
	Width     int
	Height    int
	Connected []*Room // Связанные комнаты
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
	sections := divideIntoSections(LevelWidth, LevelHeight)

	// 1. Генерируем комнаты в секциях
	for _, section := range sections {
		room := generateRoom(section)
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

func buildCorridor(r1, r2 *Room) *Corridor {
	path := []Point{}

	// Центры комнат
	x1, y1 := r1.X+r1.Width/2, r1.Y+r1.Height/2
	x2, y2 := r2.X+r2.Width/2, r2.Y+r2.Height/2

	// Определяем границы комнат (стены, где коридор может "остановиться")
	left1, right1 := r1.X, r1.X+r1.Width-1
	top1, bottom1 := r1.Y, r1.Y+r1.Height-1

	left2, right2 := r2.X, r2.X+r2.Width-1
	top2, bottom2 := r2.Y, r2.Y+r2.Height-1

	// Находим ближайшие точки входа на границе комнат
	exitX1, exitY1 := closestBoundary(x1, y1, left1, right1, top1, bottom1, x2, y2)
	exitX2, exitY2 := closestBoundary(x2, y2, left2, right2, top2, bottom2, exitX1, exitY1)

	// Горизонтально -> Вертикально
	if rand.Intn(2) == 0 {
		for x := min(exitX1, exitX2); x <= max(exitX1, exitX2); x++ {
			path = append(path, Point{X: x, Y: exitY1})
		}
		for y := min(exitY1, exitY2); y <= max(exitY1, exitY2); y++ {
			path = append(path, Point{X: exitX2, Y: y})
		}
	} else {
		// Вертикально -> Горизонтально
		for y := min(exitY1, exitY2); y <= max(exitY1, exitY2); y++ {
			path = append(path, Point{X: exitX1, Y: y})
		}
		for x := min(exitX1, exitX2); x <= max(exitX1, exitX2); x++ {
			path = append(path, Point{X: x, Y: exitY2})
		}
	}

	return &Corridor{Path: path}
}

// Вычисляет ближайшую точку на границе комнаты
func closestBoundary(cx, cy, left, right, top, bottom, targetX, targetY int) (int, int) {
	var bx, by int

	// Если цель справа, берём правый край
	if targetX > right {
		bx = right + 1
	} else if targetX < left { // Если слева, берём левый край
		bx = left - 1
	} else { // Если по горизонтали совпадает, оставляем как есть
		bx = cx
	}

	// Если цель снизу, берём нижний край
	if targetY > bottom {
		by = bottom + 1
	} else if targetY < top { // Если сверху, берём верхний край
		by = top - 1
	} else { // Если по вертикали совпадает, оставляем как есть
		by = cy
	}

	return bx, by
}

// Вспомогательные функции
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Вспомогательная функция для вычисления абсолютного значения
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}