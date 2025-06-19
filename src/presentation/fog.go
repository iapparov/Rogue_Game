package presentation

import (
	"rogue/domain"
	"math")

// ComputeFogOfWar рассчитывает видимость
func ComputeFogOfWar(player *domain.Character, level *domain.Level, fog_corr *map[string]bool) map[domain.Point]bool {
	fog := make(map[domain.Point]bool)

	for _, room := range level.Rooms {
			for y := room.Y; y < room.Y+room.Height; y++ {
				for x := room.X; x < room.X+room.Width; x++ {
					fog[domain.Point{X: x, Y: y}] = true
				}
			}
	}
	// Открываем комнату, где находится игрок

	for _, corridor := range level.Corridors {
		for _, point := range corridor.Path {
			s:=domain.Point{X: point.X, Y: point.Y}.String()
			if !(*fog_corr)[s]{
				fog[point] = true
			}
		}
	}



	for _, room := range level.Rooms {
		if player.X >= room.X && player.X < room.X+room.Width &&
			player.Y >= room.Y && player.Y < room.Y+room.Height {
			for y := room.Y; y < room.Y+room.Height; y++ {
				for x := room.X; x < room.X+room.Width; x++ {
					fog[domain.Point{X: x, Y: y}] = false
				}
			}
		}
	}
	
	// Алгоритм Брезенхэма для видимости
	for _, corridor := range level.Corridors {
		for _, point := range corridor.Path {
			if math.Abs(float64(player.X-point.X)) < 3 && math.Abs(float64(player.Y-point.Y)) < 3 {
				fog[point] = false
				s:=domain.Point{X: point.X, Y: point.Y}.String()
				(*fog_corr)[s] = true
			}
		}
	}

	for _, room := range level.Rooms {
		for _, wall := range room.Walls {
			if !fog[wall]{
				s:=domain.Point{X: wall.X, Y: wall.Y}.String()
				(*fog_corr)[s] = true
			}
		}
	}

	for _, room := range level.Rooms {
		for _, wall := range room.Walls {
			s:=domain.Point{X: wall.X, Y: wall.Y}.String()
			if (*fog_corr)[s]{
				fog[wall] = false
			}
		}
	}

	return fog
}

