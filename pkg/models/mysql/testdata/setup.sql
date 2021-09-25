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


-- Active

-- 1
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live
) VALUES (
	'Happening Now',
	'Watch now',
	'Watch the best in the world lose',
	0,
	'2021-03-07 12:00:00',
	'2321-03-17 20:00:00',
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
	is_live
) VALUES (
	'Also Now',
	'Watch now',
	'Watch the best in the world lose',
	0,
	'2021-05-07 12:00:00',
	'2321-03-17 20:00:00',
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
	is_live
) VALUES (
	'Norway Chess',
	'Watch Magnus',
	'Watch Magnus and Nepo',
	1,
	'2021-09-07 17:00:00',
	'2021-09-17 22:00:00',
	0
);

-- 4
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live
) VALUES (
	'Candidates',
	'Watch good players',
	'Watch good players fight for the chance to play against Magnus',
	1,
	'2021-03-07 11:00:00',
	'2021-03-17 12:00:00',
	0
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
	is_live
) VALUES (
	'Candidates Far Far Far',
	'Watch good players probably',
	'Watch good players probably fight for the chance to play against Magnus',
	1,
	'2321-03-07 7:00:00',
	'2321-03-17 10:00:00',
	0
);

-- 6
INSERT INTO tournaments (
	name,
	short_description,
	long_description,
	has_standings,
	start_date,
	end_date,
	is_live
) VALUES (
	'World Championship Far Far',
	'Watch Magnus vs someone',
	'Watch Magnus vs someone that might be Alireze',
	0,
	'2323-03-07 7:00:00',
	'2324-03-17 10:00:00',
	0
);
