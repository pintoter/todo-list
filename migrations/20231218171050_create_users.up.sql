ALTER TABLE notes ADD COLUMN user_id int not null;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(80) UNIQUE NOT NULL,
    login VARCHAR(80) UNIQUE NOT NULL,
    password VARCHAR(80) UNIQUE NOT NULL,
    register_at TIMESTAMP NOT NULL DEFAULT NOW(),
    refresh_token VARCHAR(255),
    expires_at TIMESTAMP,
);

ALTER TABLE notes ADD CONSTRAINT fk_notes_users FOREIGN KEY (user_id) REFERENCES users(id);