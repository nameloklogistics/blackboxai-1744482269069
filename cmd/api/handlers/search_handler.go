package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchServices handles service search requests
func (h *SearchHandler) SearchServices(c *gin.Context) {
	var filter models.SearchFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate pagination parameters
	if filter.PageSize <= 0 {
		filter.PageSize = 10 // Default page size
	}
	if filter.Page <= 0 {
		filter.Page = 1 // Default page
	}

	result, err := h.searchService.SearchServices(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSearchSuggestions handles search suggestions/autocomplete requests
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	suggestions, err := h.searchService.GetSearchSuggestions(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// RegisterRoutes registers the search routes with the given router group
func (h *SearchHandler) RegisterRoutes(router *gin.RouterGroup) {
	search := router.Group("/search")
	{
		search.POST("/services", h.SearchServices)
		search.GET("/suggestions", h.GetSearchSuggestions)
	}
}
