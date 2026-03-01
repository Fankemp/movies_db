package handlers

import (
	"filmDb/internal/repository/postgres/movies"
	"filmDb/pkg/modules"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	repo *movies.Repository
}

func NewMovieHandler(repo *movies.Repository) *MovieHandler {
	return &MovieHandler{
		repo: repo,
	}
}

func (h *MovieHandler) Create(c *gin.Context) {
	var input modules.CreateMovieRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad JSON"})
		return
	}

	err := h.repo.Save(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	moviesAll, err := h.repo.GetAllMovie(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, moviesAll)
}

func (h *MovieHandler) Search(c *gin.Context) {
	title := c.Query("title")
	ratingStr := c.Query("rating")

	if title != "" {
		movieByTitle, err := h.repo.GetByTitle(c.Request.Context(), title)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "we cant find the movie"})
			return
		}

		c.JSON(http.StatusOK, movieByTitle)
		return
	}

	if ratingStr != "" {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rating should be a number"})
			return
		}

		moviesByTitle, err := h.repo.GetByRating(c.Request.Context(), rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "some problem with DB"})
			return
		}

		c.JSON(http.StatusOK, moviesByTitle)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "write title or rating"})
}

func (h *MovieHandler) UpdateRating(c *gin.Context) {
	var rating struct {
		Title  string  `json:"title" binding:"required"`
		Rating float64 `json:"vote_average" binding:"required"`
	}

	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "write a number"})
		return
	}

	err := h.repo.UpdateRating(c.Request.Context(), rating.Title, rating.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem with DB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *MovieHandler) DeleteMovieByTitle(c *gin.Context) {
	idRow := c.Param("id")

	id, err := strconv.Atoi(idRow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should be a number"})
		return
	}

	err = h.repo.DeleteMovie(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem with DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "the movie deleted"})
}

func (h *MovieHandler) GetMovieById(c *gin.Context) {
	idRow := c.Param("id")

	id, err := strconv.Atoi(idRow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id should be a number"})
		return
	}

	movieId, err := h.repo.GetMovieById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem with DB"})
		return
	}

	c.JSON(http.StatusOK, movieId)

}
