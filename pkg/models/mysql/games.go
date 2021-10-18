package mysql

import (
	"database/sql"

	"watchess.org/watchess/pkg/models"
)

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(white, black string, res models.GameResult, whiteMatchSide, blackMatchSide, pgn string, matchID, roundID int) (int, error) {
	stmt := `INSERT INTO games (white, black, result, white_match_side, black_match_side, pgn, match_id, round_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, white, black, res.String(), whiteMatchSide, blackMatchSide, pgn, matchID, roundID)

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
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, pgn, match_id, round_id FROM games WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	g := &models.Game{}
	var res string
	err := row.Scan(
		&g.ID,
		&g.White,
		&g.Black,
		&res,
		&g.WhiteMatchSide,
		&g.BlackMatchSide,
		&g.PGN,
		&g.MatchID,
		&g.RoundID,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	gRes, err := models.GetGameResult(res)
	if err != nil {
		return nil, err
	}
	g.Result = *gRes

	return g, err
}

func (m *GameModel) GetByMatch(matchID int) ([]*models.Game, error) {
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, pgn, match_id, round_id FROM games WHERE match_id = ?`

	rows, err := m.DB.Query(stmt, matchID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []*models.Game

	for rows.Next() {
		g := &models.Game{}
		var res string
		err = rows.Scan(
			&g.ID,
			&g.White,
			&g.Black,
			&res,
			&g.WhiteMatchSide,
			&g.BlackMatchSide,
			&g.PGN,
			&g.MatchID,
			&g.RoundID,
		)
		if err != nil {
			return nil, err
		}

		gRes, err := models.GetGameResult(res)
		if err != nil {
			return nil, err
		}
		g.Result = *gRes

		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, err
}

func (m *GameModel) GetByRound(roundID int) ([]*models.Game, error) {
	stmt := `SELECT id, white, black, result, white_match_side, black_match_side, pgn, match_id, round_id FROM games WHERE round_id = ?`

	rows, err := m.DB.Query(stmt, roundID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []*models.Game

	for rows.Next() {
		g := &models.Game{}
		var res string
		err = rows.Scan(
			&g.ID,
			&g.White,
			&g.Black,
			&res,
			&g.WhiteMatchSide,
			&g.BlackMatchSide,
			&g.PGN,
			&g.MatchID,
			&g.RoundID,
		)
		if err != nil {
			return nil, err
		}

		gRes, err := models.GetGameResult(res)
		if err != nil {
			return nil, err
		}
		g.Result = *gRes

		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, err
}

func (m *GameModel) Update(id int, pgn, gameResult *string) error {
	if (pgn == nil) && (gameResult == nil) {
		return nil
	}
	stmt := "UPDATE games SET"
	var params []interface{}
	if pgn != nil {
		stmt += " pgn = ?"
		params = append(params, *pgn)
	}
	if gameResult != nil {
		if pgn != nil {
			stmt += ","
		}
		stmt += " result = ?"
		params = append(params, *gameResult)
	}
	params = append(params, id)
	stmt += " WHERE id = ?"

	result, err := m.DB.Exec(stmt, params...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if n == 0 {
		return models.ErrNoRecord
	} else if err != nil {
		return err
	}

	return nil
}

func (m *GameModel) Delete(id int) error {
	stmt := "DELETE FROM games WHERE id = ?"
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *GameModel) DeleteByRound(roundID int) error {
	stmt := "DELETE FROM games WHERE round_id = ?"
	_, err := m.DB.Exec(stmt, roundID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
