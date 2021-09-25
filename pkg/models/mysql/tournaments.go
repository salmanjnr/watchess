package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"watchess.org/watchess/pkg/models"
)

type TournamentModel struct {
	DB *sql.DB
}


// Insert row into tournaments table
func (m *TournamentModel) Insert(name, shortDescription, longDescription string, hasStandings bool, startDate, endDate time.Time, isLive bool) (int, error) {
	// Use placeholder instead of string interpolation to avoid SQL injection
	stmt := `INSERT INTO tournaments (name, short_description, long_description, has_standings,
	start_date, end_date, is_live) VALUES(?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, shortDescription, longDescription, hasStandings, startDate, endDate, isLive)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Returned id is of type int64. Convert it to int
	return int(id), nil
}

func (m *TournamentModel) Get(id int) (*models.Tournament, error) {
	stmt := `SELECT id, name, short_description, long_description, has_standings, start_date,
	end_date, is_live FROM tournaments WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	t := &models.Tournament{}

	err := row.Scan(&t.ID, &t.Name, &t.ShortDescription, &t.LongDescription, &t.HasStandings, &t.StartDate, &t.EndDate, &t.IsLive)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil{
		return nil, err
	}
	return t, nil
}


func (m *TournamentModel) scanTournaments(constraint string) ([]*models.Tournament, error) {
	stmt := "SELECT id, name, short_description, long_description, has_standings, start_date, end_date, is_live FROM tournaments " + constraint

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tournaments []*models.Tournament
	
	for rows.Next() {
		t := &models.Tournament{}
		err = rows.Scan(&t.ID, &t.Name, &t.ShortDescription, &t.LongDescription, &t.HasStandings, &t.StartDate, &t.EndDate, &t.IsLive)
		if err != nil {
			return nil, err
		}

		tournaments = append(tournaments, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
}

func (m *TournamentModel) LatestActive(limit int) ([]*models.Tournament, error) {
	stmt := fmt.Sprintf("WHERE end_date > UTC_TIMESTAMP() AND UTC_TIMESTAMP() > start_date ORDER BY start_date ASC LIMIT %d", limit)
	return m.scanTournaments(stmt)
}

func (m *TournamentModel) LatestFinished(limit int) ([]*models.Tournament, error) {
	stmt := fmt.Sprintf("WHERE end_date < UTC_TIMESTAMP() ORDER BY end_date DESC LIMIT %d", limit)

	return m.scanTournaments(stmt)
}

func (m *TournamentModel) Upcoming(limit int) ([]*models.Tournament, error) {
	stmt := fmt.Sprintf("WHERE start_date > UTC_TIMESTAMP() ORDER BY start_date ASC LIMIT %d", limit)

	return m.scanTournaments(stmt)
}
