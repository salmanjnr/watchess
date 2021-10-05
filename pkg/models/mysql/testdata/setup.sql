-- The setup could be optimized but whatever

CREATE TABLE tournaments (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(100) NOT NULL,
	short_description TEXT NOT NULL,
	long_description TEXT NOT NULL,
	has_standings BOOLEAN NOT NULL,
	start_date DATETIME NOT NULL,
	is_live BOOLEAN NOT NULL,
	end_date DATETIME NOT NULL
);

CREATE TABLE users ( 
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(50) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME NOT NULL,
	role ENUM("user", "admin") NOT NULL 
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_uc_name UNIQUE (name);

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

CREATE INDEX idx_tournament_start ON tournaments(start_date);
CREATE INDEX idx_tournament_end ON tournaments(end_date);

ALTER TABLE tournaments ADD COLUMN owner_id INTEGER NOT NULL;
ALTER TABLE tournaments ADD CONSTRAINT fk_tournament_owner_id FOREIGN KEY (owner_id) REFERENCES users(id);

CREATE INDEX idx_tournament_owner ON tournaments(owner_id);

CREATE TABLE matches (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	side1 TEXT NOT NULL,
	side2 TEXT NOT NULL,
	round_id INTEGER NOT NULL,
	CONSTRAINT fk_round_id
	FOREIGN KEY (round_id) REFERENCES rounds(id)
);

CREATE INDEX idx_match_round ON matches(round_id);

CREATE TABLE games (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	white TEXT NOT NULL,
	black TEXT NOT NULL,
	result ENUM("1-0", "0.5-0.5", "0-1"),
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

-- users

INSERT INTO users (
	name, 
	email, 
	hashed_password, 
	created, 
	role
) VALUES(
	'user', 
	'user@user.user', 
	'passpasspasspass', 
	'2000-01-01 00:01:02', 
	'admin'	
);

-- tournaments

-- Active

-- 1
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'Happening Now',
	'Watch now',
	'Watch the best in the world lose',
	0,
	'2021-03-07 12:00:00',
	'2321-03-17 20:00:00',
	1,
	1
);

-- 2
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'Also Now',
	'Watch now',
	'Watch the best in the world lose',
	0,
	'2021-05-07 12:00:00',
	'2321-03-17 20:00:00',
	1,
	1
);

-- Finished

-- 3
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'Norway Chess',
	'Watch Magnus',
	'Watch Magnus and Nepo',
	1,
	'2021-09-07 17:00:00',
	'2021-09-17 22:00:00',
	0,
	1
);

-- 4
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'Candidates',
	'Watch good players',
	'Watch good players fight for the chance to play against Magnus',
	1,
	'2021-03-07 11:00:00',
	'2021-03-17 12:00:00',
	0,
	1
);


-- Upcoming

-- 5
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'Candidates Far Far Far',
	'Watch good players probably',
	'Watch good players probably fight for the chance to play against Magnus',
	1,
	'2321-03-07 7:00:00',
	'2321-03-17 10:00:00',
	0,
	1
);

-- 6
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live,
	owner_id
) VALUES (
	'World Championship Far Far',
	'Watch Magnus vs someone',
	'Watch Magnus vs someone that might be Alireze',
	0,
	'2323-03-07 7:00:00',
	'2324-03-17 10:00:00',
	0,
	1
);

-- rounds

INSERT INTO rounds (
	name,
	pgn_source,
	white_reward_w,
	white_reward_d,
	white_reward_l,
	black_reward_w,
	black_reward_d,
	black_reward_l,
	start_date,
	tournament_id
) VALUES (
	'Round 1',
	'https://www.something.com',
	1,
	0.5,
	0,
	1,
	0.5,
	0,
	'2021-10-2 7:00:00',
	1
);

INSERT INTO rounds (
	name,
	pgn_source,
	white_reward_w,
	white_reward_d,
	white_reward_l,
	black_reward_w,
	black_reward_d,
	black_reward_l,
	start_date,
	tournament_id
) VALUES (
	'Round 10',
	'https://www.something.com',
	2,
	1,
	0,
	2,
	1,
	0,
	'2021-09-1 18:00:00',
	1
);

-- matches

INSERT INTO matches (
	side1, 
	side2, 
	round_id
) VALUES(
	'Magnus Carlsen',
	'Ian Nepo',
	1
);

INSERT INTO matches (
	side1, 
	side2, 
	round_id
) VALUES (
	'Anish Giri',
	'Levon Aronian',
	1
);

INSERT INTO matches (
	side1, 
	side2, 
	round_id
) VALUES(
	'Alireza',
	'Richard',
	2
);

-- games

INSERT INTO games (
	white,
	black,
	result,
	white_match_side,
	black_match_side,
	match_id,
	round_id
) VALUES (
	'Magnus Carlsen',
	'Ian Nepo',
	'1-0',
	'Magnus Carlsen',
	'Ian Nepo',
	1,
	1
);

INSERT INTO games (
	white,
	black,
	result,
	white_match_side,
	black_match_side,
	match_id,
	round_id
) VALUES (
	'Ian Nepo',
	'Magnus Carlsen',
	'0.5-0.5',
	'Ian Nepo',
	'Magnus Carlsen',
	1,
	1
);

INSERT INTO games (
	white,
	black,
	result,
	white_match_side,
	black_match_side,
	match_id,
	round_id
) VALUES (
	'Magnus Carlsen',
	'Ian Nepo',
	'0-1',
	'Magnus Carlsen',
	'Ian Nepo',
	1,
	1
);

INSERT INTO games (
	white,
	black,
	result,
	white_match_side,
	black_match_side,
	match_id,
	round_id
) VALUES (
	'Alireza',
	'Richard',
	'0-1',
	'Alireza',
	'Richard',
	3,
	2
);
