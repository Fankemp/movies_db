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
	title := c.Query("title")
	genre := c.Query("genre")
	orderBy := c.Query("order_by")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limits"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	var rating float64
	if ratingStr := c.Query("rating"); ratingStr != "" {
		rating, err = strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating"})
			return
		}
	}

	moviesPaginated, err := h.repo.GetPaginatedMovie(c.Request.Context(), genre, title, rating, orderBy, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, moviesPaginated)
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
func (h *MovieHandler) GetCommonRelated(c *gin.Context) {
	id1, err := strconv.Atoi(c.Query("movie_id1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie_id1"})
		return
	}

	id2, err := strconv.Atoi(c.Query("movie_id2"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie_id2"})
		return
	}

	movieCommonRelated, err := h.repo.GetCommonRelated(c.Request.Context(), id1, id2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, movieCommonRelated)
}

func (h *MovieHandler) GetDeletedMovie(c *gin.Context) {
	movieDeleted, err := h.repo.GetDeletedMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, movieDeleted)
}
