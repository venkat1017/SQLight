package storage

import (
	"encoding/json"
	"os"
	"sqlight/pkg/db"
)

// TableData represents the structure we want to save
type TableData struct {
	Tables map[string]*db.Table `json:"tables"`
}

func SaveToFile(filename string, database *db.Database) error {
	data := TableData{
		Tables: database.Tables(),
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func LoadFromFile(filename string) (*db.Database, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data TableData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	database := db.NewDatabase(filename)
	database.SetTables(data.Tables)
	return database, nil
}
