package storage

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(path string) Storage {
	db, err := sql.Open("sqlite3", filepath.Join(path, "pastey.sqlite3"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS storage (id INTEGER PRIMARY KEY, key TEXT UNIQUE NOT NULL, value TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	return &SQLiteStorage{
		db: db,
	}
}

func (c *SQLiteStorage) Save(key string, value string) {
	_, err := c.db.Exec("INSERT OR REPLACE INTO storage (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		panic(err)
	}
}

func (c *SQLiteStorage) Get(key string) (string, error) {
	row := c.db.QueryRow("SELECT value FROM storage WHERE key = ?", key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c *SQLiteStorage) Delete(key string) {
	_, err := c.db.Exec("DELETE FROM storage WHERE key = ?", key)
	if err != nil {
		panic(err)
	}
}

func (c *SQLiteStorage) Exists(key string) bool {
	row := c.db.QueryRow("SELECT key FROM storage WHERE key = ?", key)

	var k string
	err := row.Scan(&k)

	return err == nil
}

func (c *SQLiteStorage) Close() {
	c.db.Close()
}
