CREATE TYPE task_status AS ENUM ('done', 'not done');

CREATE TABLE IF NOT EXISTS notes (
  id SERIAL PRIMARY KEY,
  title VARCHAR(80) NOT NULL,
  description TEXT DEFAULT NULL, 
  date TIMESTAMP NOT NULL,
  status task_status DEFAULT 'not done'
);