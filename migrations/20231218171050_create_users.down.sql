ALTER TABLE IF EXISTS notes DROP CONSTRAINT IF EXISTS fk_notes_users;

DROP TABLE IF EXISTS users;

ALTER TABLE notes DROP COLUMN IF EXISTS user_id;