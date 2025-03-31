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

func Json_Load() (*domain.GameSession, error){
	file, err := os.Open("Save.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var session domain.GameSession
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&session)
	return &session, err
}

func Json_Free() error {
	err := os.WriteFile("Save.json", []byte("[]"), 0644) // Для массива
	return err
}

// Загружает данные рекордов из файла
func Json_Record_Load(scoreBoard *domain.LeaderBoards) error {
	// Открываем файл, если он существует
	file, err := os.Open("Records.json")
	if os.IsNotExist(err) { // Если файла нет, создаем новый
		file, err = os.Create("Records.json")
		if err != nil {
			return err
		}
		defer file.Close()

		// Записываем пустой JSON в файл
		emptyData := domain.LeaderBoards{Record: []*domain.Records{}}
		jsonData, _ := json.MarshalIndent(emptyData, "", "  ")
		file.Write(jsonData) // Записываем пустой массив записей

		*scoreBoard = emptyData // Устанавливаем пустую структуру
		return nil
	} else if err != nil {
		return err
	}
	defer file.Close()

	// Декодируем JSON в структуру
	decoder := json.NewDecoder(file)
	return decoder.Decode(scoreBoard)
}

func Json_Record_Save(game *domain.GameSession, ScoreBoard *domain.LeaderBoards) error {
	file, err := os.Create("Records.json")
	if err != nil {
		return err
	}
	defer file.Close()
	newRecord := domain.Records{Name: game.Player.Name, Record: game.TreasureCount}
	ScoreBoard.Record = append(ScoreBoard.Record, &newRecord)
	encoder:= json.NewEncoder(file)
	encoder.SetIndent(""," ")
	return encoder.Encode(ScoreBoard)
}