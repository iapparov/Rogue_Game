package domain

import "math/rand"

// Attack вычисляет атаку
func (c *Character) Attack(target *Character) bool {
	if rand.Float64() < calculateHitChance(c, target) {
		damage := calculateDamage(c)
		target.Health -= damage
		return true
	}
	return false
}

// calculateHitChance рассчитывает вероятность попадания
func calculateHitChance(attacker, defender *Character) float64 {
	return 0.5 + float64(attacker.Agility-defender.Agility)*0.05
}

// calculateDamage рассчитывает урон
func calculateDamage(attacker *Character) int {
	baseDamage := attacker.Strength
	if attacker.Weapon != nil {
		baseDamage += attacker.Weapon.Strength
	}
	return baseDamage
}