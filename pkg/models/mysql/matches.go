package mysql

import (
	"database/sql"
	"log"

	"watchess.org/watchess/pkg/models"
)

type MatchModel struct {
	DB *sql.DB
}

func (m *MatchModel) Insert(side1, side2 string, roundID int) (int, error) {
	stmt := `INSERT INTO matches (side1, side2, round_id) VALUES(?, ?, ?)`

	log.Printf("Round ID: %d\n", roundID)
	result, err := m.DB.Exec(stmt, side1, side2, roundID)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *MatchModel) Get(id int) (*models.Match, error) {
	stmt := `SELECT id, side1, side2, round_id FROM matches WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	c := &models.Match{}

	err := row.Scan(&c.ID, &c.Side1, &c.Side2, &c.RoundID)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return c, nil
}
