package modules

import "time"

type Movie struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"original_title" json:"title"`
	Language    string     `db:"original_language" json:"language"`
	Overview    string     `db:"overview" json:"overview"`
	Genre       string     `db:"genre" json:"genre"`
	ReleaseDate time.Time  `db:"release_date" json:"release_date"`
	Rating      float64    `db:"vote_average" json:"rating"`
	VoteCount   int        `db:"vote_count" json:"vote_count"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}
type CreateMovieRequest struct {
	Title       string    `json:"title" binding:"required"`
	Language    string    `json:"language"`
	Overview    string    `json:"overview"`
	ReleaseDate time.Time `json:"release_date"`
}
