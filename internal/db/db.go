package db

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ChatInfo struct {
	ChatID int64
}

type DB struct {
	path     string
	chatsIDs map[int64]ChatInfo
}

// NewDB open and load db if exist, or create new one if not exist
func NewDB(path string) (DB, error) {
	if err := os.MkdirAll(path, 0777); err != nil {
		return DB{}, err
	}

	chatsIDs := make(map[int64]ChatInfo)
	path = filepath.Join(path, "db.json")

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) || len(data) == 0 {
		return DB{
			path:     path,
			chatsIDs: chatsIDs,
		}, nil
	}
	if err != nil {
		return DB{}, err
	}

	if err := json.Unmarshal(data, &chatsIDs); err != nil {
		return DB{}, err
	}

	return DB{
		path:     path,
		chatsIDs: chatsIDs,
	}, nil
}

func (db *DB) Save() error {
	file, err := json.MarshalIndent(db.chatsIDs, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, file, 0777)
}

func (db *DB) Get(key int64) ChatInfo {
	return db.chatsIDs[key]
}

func (db *DB) Update(key int64, value ChatInfo) {
	db.chatsIDs[key] = value
}

func (db *DB) Delete(key int64) {
	delete(db.chatsIDs, key)
}

func (db *DB) List() []ChatInfo {
	chatsInfo := make([]ChatInfo, 0, len(db.chatsIDs))
	for _, val := range db.chatsIDs {
		chatsInfo = append(chatsInfo, val)
	}
	return chatsInfo
}
