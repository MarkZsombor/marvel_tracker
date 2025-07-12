package models

import (
	"database/sql"
	"time"
)

type Hero struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Scenario struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Play struct {
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	Outcome    string    `json:"outcome"`
	Difficulty string    `json:"difficulty"`
	Notes      string    `json:"notes"`
	ScenarioID int       `json:"scenario_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Deck struct {
	ID        int       `json:"id"`
	PlayID    int       `json:"play_id"`
	HeroID    int       `json:"hero_id"`
	Aspect    string    `json:"aspect"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlayRepository struct {
	db *sql.DB
}

func NewPlayRepository(db *sql.DB) *PlayRepository {
	return &PlayRepository{db: db}
}

func (r *PlayRepository) GetAll() ([]Play, error) {
	rows, err := r.db.Query("SELECT id, date, outcome, difficulty, notes, scenario_id, created_at, updated_at FROM plays ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plays []Play
	for rows.Next() {
		var p Play
		err := rows.Scan(&p.ID, &p.Date, &p.Outcome, &p.Difficulty, &p.Notes, &p.ScenarioID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		plays = append(plays, p)
	}

	return plays, nil
}

func (r *PlayRepository) Create(p *Play) error {
	result, err := r.db.Exec(
		"INSERT INTO plays (date, outcome, difficulty, notes, scenario_id) VALUES (?, ?, ?, ?, ?)",
		p.Date, p.Outcome, p.Difficulty, p.Notes, p.ScenarioID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}