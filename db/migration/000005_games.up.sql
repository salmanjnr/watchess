CREATE TABLE games (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	white TEXT NOT NULL,
	black TEXT NOT NULL,
	result ENUM("*", "1-0", "1/2-1/2", "0-1") DEFAULT "*" NOT NULL,
	white_match_side TEXT NOT NULL,
	black_match_side TEXT NOT NULL,
	match_id INTEGER NOT NULL,
	CONSTRAINT fk_game_match_id
	FOREIGN KEY (match_id) REFERENCES matches(id),
	round_id INTEGER NOT NULL,
	CONSTRAINT fk_game_round_id
	FOREIGN KEY (round_id) REFERENCES rounds(id)
);

CREATE INDEX idx_game_match ON games(match_id);
CREATE INDEX idx_game_round ON games(round_id);
