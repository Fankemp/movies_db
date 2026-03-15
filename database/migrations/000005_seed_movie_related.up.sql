INSERT INTO movie_related (movie_id, related_id)
SELECT m1.id, m2.id
FROM movies m1
JOIN movies m2 ON m1.genre = m2.genre
WHERE m1.id <> m2.id
ON CONFLICT DO NOTHING
