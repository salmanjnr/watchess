package mysql

import (
	"database/sql"

	"watchess.org/watchess/pkg/models"
)

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(white, black string, res *models.GameResult, whiteMatchSide, blackMatchSide string, matchID, roundID int) (int, error) {
	stmt := `INSERT INTO games (white, black, result, white_match_side, black_match_side, match_id, round_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	var resStr sql.NullString

	if res != nil {
		resStr = sql.NullString{String: res.String(), Valid: true}
	} else {
		resStr = sql.NullString{}
	}

	result, err := m.DB.Exec(stmt, white, black, resStr, whiteMatchSide, blackMatchSide, matchID, roundID)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *GameModel) Get(id int) (*models.Game, error) {
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, match_id, round_id FROM games WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	g := &models.Game{}
	var res sql.NullString
	err := row.Scan(
		&g.ID,
		&g.White,
		&g.Black,
		&res,
		&g.WhiteMatchSide,
		&g.BlackMatchSide,
		&g.MatchID,
		&g.RoundID,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	
	if res.Valid {
		g.Result, err = models.GetGameResult(res.String)
		if err != nil {
			return nil ,err
		}
	}

	return g, err
}

func (m *GameModel) GetByMatch(matchID int) ([]*models.Game, error) {
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, match_id, round_id FROM games WHERE match_id = ?`

	rows, err := m.DB.Query(stmt, matchID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []*models.Game

	for rows.Next() {
		g := &models.Game{}
		var res sql.NullString
		err = rows.Scan(
			&g.ID,
			&g.White,
			&g.Black,
			&res,
			&g.WhiteMatchSide,
			&g.BlackMatchSide,
			&g.MatchID,
			&g.RoundID,
		)
		if err != nil {
			return nil, err
		}

		if res.Valid {
			g.Result, err = models.GetGameResult(res.String)
			if err != nil {
				return nil ,err
			}
		}

		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, err
}

func (m *GameModel) GetByRound(roundID int) ([]*models.Game, error) {
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, match_id, round_id FROM games WHERE round_id = ?`

	rows, err := m.DB.Query(stmt, roundID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []*models.Game

	for rows.Next() {
		g := &models.Game{}
		var res sql.NullString
		err = rows.Scan(
			&g.ID,
			&g.White,
			&g.Black,
			&res,
			&g.WhiteMatchSide,
			&g.BlackMatchSide,
			&g.MatchID,
			&g.RoundID,
		)
		if err != nil {
			return nil, err
		}

		if res.Valid {
			g.Result, err = models.GetGameResult(res.String)
			if err != nil {
				return nil ,err
			}
		}

		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, err
}
