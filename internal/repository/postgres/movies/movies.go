package movies

import (
	"context"
	"filmDb/internal/repository/postgres"
	"filmDb/pkg/modules"
	"fmt"
	"time"
)

type Repository struct {
	db               *postgres.Storage
	executionTimeout time.Duration
}

func NewRepository(db *postgres.Storage) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) Save(ctx context.Context, movie *modules.CreateMovieRequest) error {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	query := `INSERT INTO movies (original_title, original_language, overview, release_date, vote_average, vote_count)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.DB.ExecContext(ctx, query,
		movie.Title, movie.Language, movie.Overview, movie.ReleaseDate, 0.0, 0)

	return err
}

func (r *Repository) DeleteMovie(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	query := `UPDATE movies SET deleted_at = NOW() WHERE id = $1`
	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error with delete: %w", err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("we cant get the deleted rows: %w", err)
	}

	if row == 0 {
		return fmt.Errorf("movie with id %d not found", id)
	}

	return nil
}

func (r *Repository) UpdateRating(ctx context.Context, title string, rating float64) error {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	query := `UPDATE movies SET vote_average = $1 WHERE original_title = $2`
	_, err := r.db.DB.ExecContext(ctx, query, rating, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMovieById(ctx context.Context, id int) (modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movie modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count
				FROM movies
				WHERE id = $1`

	err := r.db.DB.GetContext(ctx, &movie, query, id)
	if err != nil {
		return modules.Movie{}, err
	}

	return movie, nil
}

func (r *Repository) GetDeletedMovies(ctx context.Context) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	movieDeleted := make([]modules.Movie, 0)

	query := `SELECT * FROM movies WHERE deleted_at IS NOT NULL`
	err := r.db.DB.SelectContext(ctx, &movieDeleted, query)
	if err != nil {
		return nil, err
	}

	return movieDeleted, nil
}

func (r *Repository) GetPaginatedMovie(
	ctx context.Context,
	genre string,
	title string,
	rating float64,
	orderBy string,
	limit int,
	offset int,
) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	moviePaginated := make([]modules.Movie, 0)
	args := make([]interface{}, 0)
	query := `SELECT * FROM movies WHERE 1=1 AND deleted_at IS NULL`

	if genre != "" {
		args = append(args, genre)
		query += fmt.Sprintf(` AND genre = $%d`, len(args))
	}

	if title != "" {
		args = append(args, title)
		query += fmt.Sprintf(` AND original_title ILIKE $%d`, len(args))
	}

	if rating > 0 {
		args = append(args, rating)
		query += fmt.Sprintf(`AND vote_average >= $%d`, len(args))
	}

	allowedOrderBy := map[string]bool{
		"id":           true,
		"vote_average": true,
		"release_data": true,
		"vote_count":   true,
	}

	if !allowedOrderBy[orderBy] {
		orderBy = "vote_average"
	}
	query += fmt.Sprintf(` ORDER BY %s`, orderBy)

	args = append(args, limit, offset)
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)-1, len(args))

	err := r.db.DB.SelectContext(ctx, &moviePaginated, query, args...)
	if err != nil {
		return nil, err
	}

	return moviePaginated, nil
}

func (r *Repository) GetCommonRelated(ctx context.Context, movieID1, movieID2 int) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()
	moviesRelated := make([]modules.Movie, 0)

	query := `SELECT m.*
			  FROM movies m
			  JOIN movie_related mr1 ON m.id = mr1.related_id
			  JOIN movie_related mr2 ON m.id = mr2.related_id 
			  WHERE mr1.movie_id = $1 AND mr2.movie_id = $2
			  `
	err := r.db.DB.SelectContext(ctx, &moviesRelated, query, movieID1, movieID2)
	if err != nil {
		return nil, err
	}

	return moviesRelated, nil
}
