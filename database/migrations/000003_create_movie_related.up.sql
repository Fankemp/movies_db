CREATE TABLE IF NOT EXISTS movie_related(
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    related_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, related_id),
    CONSTRAINT check_not_sel_related CHECK ( movie_id <> related_id)
);