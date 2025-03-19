package domain

type EnemyType string
const (
	Zombie    EnemyType = "Zombie"
	Vampire   EnemyType = "Vampire"
	Ghost     EnemyType = "Ghost"
	Ogre      EnemyType = "Ogre"
	SnakeMage EnemyType = "SnakeMage"
)


type Enemy struct{
	Character
	Type EnemyType
	Hostility int //враждебность
	X, Y int
}

// NewEnemy создаёт нового врага
func NewEnemy(enemyType EnemyType) *Enemy {
	health, agility, strength, hostility := getEnemyStats(enemyType)

	return &Enemy{
		Character: Character{
			Name:      string(enemyType),
			Health:    health,
			MaxHealth: health,
			Agility: agility,
			Strength:  strength,
		},
		Type:      enemyType,
		Hostility: hostility,
	}
}

// getEnemyStats возвращает характеристики противника
func getEnemyStats(enemyType EnemyType) (int, int, int, int) {
	switch enemyType {
	case Zombie:
		return 50, 2, 5, 3
	case Vampire:
		return 40, 8, 4, 9
	case Ghost:
		return 20, 10, 3, 5
	case Ogre:
		return 80, 3, 12, 6
	case SnakeMage:
		return 30, 9, 6, 10
	default:
		return 30, 5, 5, 5
	}
}

