import { api } from './api';
import { SearchFilter, SearchResult, SearchSuggestion } from '../types/search';
import { ApiResponse } from '../types';

export const searchApi = {
  searchServices: async (filter: SearchFilter): Promise<ApiResponse<SearchResult>> => {
    const response = await api.client.get('/api/search/services', { params: filter });
    return response.data;
  },

  getSearchSuggestions: async (query: string): Promise<ApiResponse<SearchSuggestion[]>> => {
    const response = await api.client.get('/api/search/suggestions', {
      params: { query },
    });
    return response.data;
  },
};

export default searchApi;
