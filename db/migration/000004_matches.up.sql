CREATE TABLE matches (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	side1 TEXT NOT NULL,
	side2 TEXT NOT NULL,
	round_id INTEGER NOT NULL,
	CONSTRAINT fk_round_id
	FOREIGN KEY (round_id) REFERENCES rounds(id)
);

CREATE INDEX idx_match_round ON matches(round_id);
