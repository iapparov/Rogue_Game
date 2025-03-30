package datalayer

import (
	"encoding/json"
	"os"
	"rogue/domain"
)

func Json_Save(game *domain.GameSession) error {
	file, err := os.Create("Save.json")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")
	return encoder.Encode(game)

}
