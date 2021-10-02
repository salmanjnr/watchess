ALTER TABLE tournaments ADD COLUMN owner_id INTEGER NOT NULL;
ALTER TABLE tournaments ADD CONSTRAINT fk_tournament_owner_id FOREIGN KEY (owner_id) REFERENCES users(id);

CREATE INDEX idx_tournament_owner ON tournaments(owner_id);
