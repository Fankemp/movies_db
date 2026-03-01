CREATE TABLE if not exists movies (
        id SERIAL PRIMARY KEY,
        original_title VARCHAR(255) NOT NULL,
        original_language VARCHAR(10),
        overview TEXT,
        release_date DATE,
        vote_average DECIMAL(3, 1),
        vote_count INT
);

CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(original_title);
