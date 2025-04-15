package services

import (
	"fmt"
	"strings"
	"time"
	
	"logistics-marketplace/internal/models"
)

type SearchService struct {
	marketplaceService *MarketplaceService
}

func NewSearchService(marketplaceService *MarketplaceService) *SearchService {
	return &SearchService{
		marketplaceService: marketplaceService,
	}
}

// SearchServices searches for logistics services based on the provided filters
func (s *SearchService) SearchServices(filter *models.SearchFilter) (*models.SearchResult, error) {
	if filter.PageSize == 0 {
		filter.PageSize = 10 // Default page size
	}
	if filter.Page == 0 {
		filter.Page = 1 // Default page
	}

	// Get all services first (in a real implementation, this would be a database query)
	services, err := s.marketplaceService.GetServicesByCategory(string(filter.Category))
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	// Apply filters
	filtered := s.applyFilters(services, filter)

	// Apply sorting
	sorted := s.applySorting(filtered, filter)

	// Apply pagination
	paginated, totalCount := s.applyPagination(sorted, filter)

	// Calculate total pages
	totalPages := (totalCount + filter.PageSize - 1) / filter.PageSize

	return &models.SearchResult{
		Services:    paginated,
		TotalCount:  totalCount,
		CurrentPage: filter.Page,
		PageSize:    filter.PageSize,
		TotalPages:  totalPages,
	}, nil
}

// GetSearchSuggestions returns search suggestions based on the input
func (s *SearchService) GetSearchSuggestions(query string) ([]models.SearchSuggestion, error) {
	suggestions := make([]models.SearchSuggestion, 0)

	// Add location suggestions
	locations := s.searchLocations(query)
	for _, loc := range locations {
		suggestions = append(suggestions, models.SearchSuggestion{
			Type:  "location",
			Value: loc,
			Label: fmt.Sprintf("Location: %s", loc),
		})
	}

	// Add service suggestions
	services := s.searchServiceTypes(query)
	for _, svc := range services {
		suggestions = append(suggestions, models.SearchSuggestion{
			Type:  "service",
			Value: svc,
			Label: fmt.Sprintf("Service: %s", svc),
		})
	}

	return suggestions, nil
}

// Helper functions

func (s *SearchService) applyFilters(services []models.LogisticsService, filter *models.SearchFilter) []models.LogisticsService {
	filtered := make([]models.LogisticsService, 0)

	for _, service := range services {
		if !s.matchesFilter(service, filter) {
			continue
		}
		filtered = append(filtered, service)
	}

	return filtered
}

func (s *SearchService) matchesFilter(service models.LogisticsService, filter *models.SearchFilter) bool {
	// Check category
	if filter.Category != "" && service.Category != filter.Category {
		return false
	}

	// Check subcategory
	if filter.SubCategory != "" && service.SubCategory.ID != filter.SubCategory {
		return false
	}

	// Check origin/destination
	if filter.Origin != "" && !strings.Contains(strings.ToLower(service.Origin), strings.ToLower(filter.Origin)) {
		return false
	}
	if filter.Destination != "" && !strings.Contains(strings.ToLower(service.Destination), strings.ToLower(filter.Destination)) {
		return false
	}

	// Check price range
	if filter.MinPrice != nil && service.Item.BasePrice < *filter.MinPrice {
		return false
	}
	if filter.MaxPrice != nil && service.Item.BasePrice > *filter.MaxPrice {
		return false
	}

	// Check route type
	if filter.RouteType != "" && service.RouteType != filter.RouteType {
		return false
	}

	// Check provider
	if filter.ProviderID != "" && service.Provider.ID != filter.ProviderID {
		return false
	}

	// Check active status
	if filter.IsActive != nil && service.IsActive != *filter.IsActive {
		return false
	}

	// Check validity period
	if !filter.ValidFrom.IsZero() && service.Item.ValidFrom.Before(filter.ValidFrom) {
		return false
	}
	if !filter.ValidUntil.IsZero() && service.Item.ValidUntil.After(filter.ValidUntil) {
		return false
	}

	// Check general query
	if filter.Query != "" {
		query := strings.ToLower(filter.Query)
		if !strings.Contains(strings.ToLower(service.Item.Name), query) &&
			!strings.Contains(strings.ToLower(service.Item.Description), query) &&
			!strings.Contains(strings.ToLower(service.Provider.Name), query) {
			return false
		}
	}

	return true
}

func (s *SearchService) applySorting(services []models.LogisticsService, filter *models.SearchFilter) []models.LogisticsService {
	if filter.SortBy == "" {
		return services
	}

	sorted := make([]models.LogisticsService, len(services))
	copy(sorted, services)

	// Sort based on the specified field
	switch filter.SortBy {
	case "price":
		if filter.SortOrder == "desc" {
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Item.BasePrice > sorted[j].Item.BasePrice
			})
		} else {
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Item.BasePrice < sorted[j].Item.BasePrice
			})
		}
	case "transit_time":
		if filter.SortOrder == "desc" {
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].TransitTime > sorted[j].TransitTime
			})
		} else {
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].TransitTime < sorted[j].TransitTime
			})
		}
	}

	return sorted
}

func (s *SearchService) applyPagination(services []models.LogisticsService, filter *models.SearchFilter) ([]models.LogisticsService, int) {
	totalCount := len(services)

	start := (filter.Page - 1) * filter.PageSize
	if start >= totalCount {
		return []models.LogisticsService{}, totalCount
	}

	end := start + filter.PageSize
	if end > totalCount {
		end = totalCount
	}

	return services[start:end], totalCount
}

func (s *SearchService) searchLocations(query string) []string {
	// In a real implementation, this would query a location database or API
	// This is a placeholder implementation
	commonLocations := []string{
		"Singapore",
		"Hong Kong",
		"Shanghai",
		"Tokyo",
		"Seoul",
		"Bangkok",
		"Jakarta",
		"Mumbai",
		"Dubai",
		"Rotterdam",
	}

	if query == "" {
		return commonLocations
	}

	matches := make([]string, 0)
	query = strings.ToLower(query)
	for _, loc := range commonLocations {
		if strings.Contains(strings.ToLower(loc), query) {
			matches = append(matches, loc)
		}
	}

	return matches
}

func (s *SearchService) searchServiceTypes(query string) []string {
	// In a real implementation, this would query a service type database or API
	// This is a placeholder implementation
	serviceTypes := []string{
		"FCL Import",
		"LCL Import",
		"Air Import",
		"FCL Export",
		"LCL Export",
		"Air Export",
		"Land Transit",
		"Rail Transit",
		"Sea-Sea Transshipment",
		"Air-Air Transshipment",
	}

	if query == "" {
		return serviceTypes
	}

	matches := make([]string, 0)
	query = strings.ToLower(query)
	for _, svc := range serviceTypes {
		if strings.Contains(strings.ToLower(svc), query) {
			matches = append(matches, svc)
		}
	}

	return matches
}
