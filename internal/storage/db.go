package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        login TEXT UNIQUE,
        password TEXT
    )`); err != nil {
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS expressions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        expression TEXT,
        result REAL,
        error TEXT,
        FOREIGN KEY(user_id) REFERENCES users(id)
    )`); err != nil {
		return nil, fmt.Errorf("failed to create expressions table: %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) RegisterUser(login, password string) error {
	_, err := s.db.Exec("INSERT INTO users (login, password) VALUES (?, ?)", login, password)
	if err != nil {
		return errors.New("failed to register user")
	}
	return nil
}

func (s *Storage) GetUserByLogin(login string) (*User, error) {
	var u User
	row := s.db.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login)
	if err := row.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) SaveExpression(userID int, expr string, result float64, errStr string) error {
	_, err := s.db.Exec(
		"INSERT INTO expressions (user_id, expression, result, error) VALUES (?, ?, ?, ?)",
		userID, expr, result, errStr,
	)
	return err
}

func (s *Storage) ListExpressions(userID int) ([]*Expression, error) {
	rows, err := s.db.Query(
		"SELECT id, expression, result, error FROM expressions WHERE user_id = ? ORDER BY id",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Expression
	for rows.Next() {
		e := &Expression{UserID: userID}
		if err := rows.Scan(&e.ID, &e.Expression, &e.Result, &e.Error); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (s *Storage) GetExpression(userID, exprID int) (*Expression, error) {
	row := s.db.QueryRow(
		"SELECT id, expression, result, error FROM expressions WHERE user_id = ? AND id = ?",
		userID, exprID,
	)

	var e Expression
	e.UserID = userID
	if err := row.Scan(&e.ID, &e.Expression, &e.Result, &e.Error); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}
