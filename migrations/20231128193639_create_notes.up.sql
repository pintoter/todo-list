CREATE TYPE note_status AS ENUM ('done', 'not_done');

CREATE TABLE IF NOT EXISTS notes (
  id SERIAL PRIMARY KEY,
  title VARCHAR(80) NOT NULL,
  description TEXT DEFAULT NULL, 
  date DATE NOT NULL DEFAULT CURRENT_DATE,
  status note_status DEFAULT 'not_done'
);