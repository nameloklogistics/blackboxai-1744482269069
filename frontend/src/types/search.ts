
import { ServiceType, ShipmentMode, ServiceListing } from './index';
export interface SearchFilter {
  query?: string;
  serviceType?: ServiceType;
  shipmentMode?: ShipmentMode;
  origin?: string;
  destination?: string;
  minPrice?: number;
  maxPrice?: number;
  routeType?: string;
  transitTime?: string;
  providerId?: string;
  isActive?: boolean;
  validFrom?: string;
  validUntil?: string;
  page: number;
  pageSize: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface SearchResult {
  services: ServiceListing[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export interface SearchSuggestion {
  type: 'location' | 'service' | 'provider';
  value: string;
  label: string;
}

export interface SearchStats {
  popularSearches: Array<{
    query: string;
    count: number;
    lastUsed: string;
  }>;
  popularFilters: Array<{
    filter: string;
    value: string;
    count: number;
  }>;
  popularRoutes: Array<{
    origin: string;
    destination: string;
    count: number;
  }>;
}

export interface SearchState {
  query: string;
  filters: SearchFilter;
  results: SearchResult | null;
  suggestions: SearchSuggestion[];
  loading: boolean;
  error: string | null;
}
