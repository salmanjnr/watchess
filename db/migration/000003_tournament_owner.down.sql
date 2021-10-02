ALTER TABLE tournaments
DROP FOREIGN KEY fk_tournament_owner_id;
ALTER TABLE tournaments
DROP COLUMN IF EXISTS owner_id;
