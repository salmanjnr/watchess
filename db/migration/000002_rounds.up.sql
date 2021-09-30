CREATE TABLE rounds (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(10) NOT NULL,
	pgn_source TEXT NOT NULL,
	white_reward_w FLOAT NOT NULL,
	white_reward_d FLOAT NOT NULL,
	white_reward_l FLOAT NOT NULL,
	black_reward_w FLOAT NOT NULL,
	black_reward_d FLOAT NOT NULL,
	black_reward_l FLOAT NOT NULL,
	start_date DATETIME NOT NULL,
	tournament_id INTEGER NOT NULL,
	FOREIGN KEY (tournament_id) REFERENCES tournaments(id)
);

CREATE INDEX idx_round_tournament ON rounds(tournament_id);

-- Also create indices on tournaments dates
CREATE INDEX idx_tournament_start ON tournaments(start_date);
CREATE INDEX idx_tournament_end ON tournaments(end_date);
