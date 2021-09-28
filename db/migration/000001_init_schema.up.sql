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
