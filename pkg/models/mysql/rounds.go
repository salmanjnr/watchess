package mysql

import (
	"database/sql"
	"time"

	"watchess.org/watchess/pkg/models"
)

type RoundModel struct {
	DB *sql.DB
}

func (m *RoundModel) Insert(name, pgn_source string, whiteReward, blackReward models.GameReward, startDate time.Time, tournament_id int) (int, error) {
	stmt := `INSERT INTO rounds (name, pgn_source, white_reward_w, white_reward_d, white_reward_l,
	black_reward_w, black_reward_d, black_reward_l, start_date, tournament_id)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, name, pgn_source, whiteReward.Win, whiteReward.Draw, whiteReward.Loss, blackReward.Win, blackReward.Draw, blackReward.Loss, startDate, tournament_id)

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

func (m *RoundModel) Get(id int) (*models.Round, error) {
	stmt := `SELECT id, name, pgn_source, white_reward_w, white_reward_d, white_reward_l,
	black_reward_w, black_reward_d, black_reward_l,
	start_date, tournament_id FROM rounds WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	r := &models.Round{}

	err := row.Scan(&r.ID, &r.Name, &r.PGNSource, &(r.WhiteReward.Win), &(r.WhiteReward.Draw), &(r.WhiteReward.Loss), &(r.BlackReward.Win), &(r.BlackReward.Draw), &(r.BlackReward.Loss), &r.StartDate, &r.TournamentID)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *RoundModel) GetByTournament(tournamentId int) ([]*models.Round, error) {
	stmt := `SELECT id, name, pgn_source, white_reward_w, white_reward_d, white_reward_l,
	black_reward_w, black_reward_d, black_reward_l,
	start_date, tournament_id FROM rounds WHERE tournament_id = ?`
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rounds []*models.Round

	for rows.Next() {
		r := &models.Round{}
		err = rows.Scan(&r.ID, &r.Name, &r.PGNSource, &(r.WhiteReward.Win), &(r.WhiteReward.Draw), &(r.WhiteReward.Loss), &(r.BlackReward.Win), &(r.BlackReward.Draw), &(r.BlackReward.Loss), &r.TournamentID)
		if err != nil {
			return nil, err
		}

		rounds = append(rounds, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rounds, nil
}
