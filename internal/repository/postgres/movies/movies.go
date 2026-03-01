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

	query := `DELETE FROM movies WHERE id = $1`
	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error with delete: %w", err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("we cant get the deleted rows: %w", err)
	}

	if row == 0 {
		return fmt.Errorf("there is not movie in DB: %w", err)
	}

	return nil
}

func (r *Repository) GetAllMovie(ctx context.Context) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movies []modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count
			   FROM movies`
	err := r.db.DB.SelectContext(ctx, &movies, query)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) GetByTitle(ctx context.Context, title string) (*modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movie modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count 
          FROM movies 
          WHERE original_title ILIKE $1 
          LIMIT 1`
	err := r.db.DB.GetContext(ctx, &movie, query, title)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *Repository) GetByDate(ctx context.Context, release time.Time) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movie []modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count 
          FROM movies 
          WHERE release_date  = $1::date`

	err := r.db.DB.SelectContext(ctx, &movie, query, release)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *Repository) GetByLanguage(ctx context.Context, language string) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movies []modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count 
			  FROM movies
			  WHERE original_language = $1`
	err := r.db.DB.SelectContext(ctx, &movies, query, language)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) GetByRating(ctx context.Context, rating float64) ([]modules.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var movies []modules.Movie

	query := `SELECT id, original_title, original_language, 
          overview, release_date, vote_average, vote_count
				FROM movies
				WHERE vote_average = $1`

	err := r.db.DB.SelectContext(ctx, &movies, query, rating)
	if err != nil {
		return nil, err
	}

	return movies, nil
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
