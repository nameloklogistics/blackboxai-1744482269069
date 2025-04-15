import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ServiceListing, Booking, ApiResponse } from '../../types';
import { SearchFilter, SearchResult, SearchState } from '../../types/search';
import { api } from '../../services/api';
import { searchApi } from '../../services/search';

interface MarketplaceState extends SearchState {
  listings: ServiceListing[];
  bookings: Booking[];
  selectedListing: ServiceListing | null;
  selectedBooking: Booking | null;
}

const initialState: MarketplaceState = {
  listings: [],
  bookings: [],
  selectedListing: null,
  selectedBooking: null,
  loading: false,
  error: null,
  query: '',
  filters: {
    page: 1,
    pageSize: 10,
  },
  results: null,
  suggestions: [],
};

// Error handler helper
const handleError = (error: any, rejectWithValue: Function, defaultMessage: string) => {
  if (error?.response?.data?.message) {
    return rejectWithValue(error.response.data.message);
  }
  return rejectWithValue(defaultMessage);
};

// Async thunks
export const searchServices = createAsyncThunk<
  SearchResult,
  SearchFilter,
  { rejectValue: string }
>('marketplace/searchServices', async (filter, { rejectWithValue }) => {
  try {
    const response = await searchApi.searchServices(filter);
    if (response.success) {
      return response.data;
    }
    return rejectWithValue(response.message || 'Failed to search services');
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to search services');
  }
});

export const fetchServiceListings = createAsyncThunk<
  ApiResponse<ServiceListing[]>,
  void,
  { rejectValue: string }
>('marketplace/fetchServiceListings', async (_, { rejectWithValue }) => {
  try {
    return await api.getServiceListings();
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to fetch service listings');
  }
});

export const createServiceListing = createAsyncThunk<
  ApiResponse<ServiceListing>,
  Partial<ServiceListing>,
  { rejectValue: string }
>('marketplace/createServiceListing', async (listing, { rejectWithValue }) => {
  try {
    return await api.createServiceListing(listing);
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to create service listing');
  }
});

export const updateServiceListing = createAsyncThunk<
  ApiResponse<ServiceListing>,
  { id: string; updates: Partial<ServiceListing> },
  { rejectValue: string }
>('marketplace/updateServiceListing', async ({ id, updates }, { rejectWithValue }) => {
  try {
    return await api.updateServiceListing(id, updates);
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to update service listing');
  }
});

export const deleteServiceListing = createAsyncThunk<
  void,
  string,
  { rejectValue: string }
>('marketplace/deleteServiceListing', async (id, { rejectWithValue }) => {
  try {
    await api.deleteServiceListing(id);
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to delete service listing');
  }
});

export const fetchBookings = createAsyncThunk<
  ApiResponse<Booking[]>,
  void,
  { rejectValue: string }
>('marketplace/fetchBookings', async (_, { rejectWithValue }) => {
  try {
    return await api.getBookings();
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to fetch bookings');
  }
});

export const createBooking = createAsyncThunk<
  ApiResponse<Booking>,
  Partial<Booking>,
  { rejectValue: string }
>('marketplace/createBooking', async (booking, { rejectWithValue }) => {
  try {
    return await api.createBooking(booking);
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to create booking');
  }
});

export const updateBookingStatus = createAsyncThunk<
  ApiResponse<Booking>,
  { id: string; status: string },
  { rejectValue: string }
>('marketplace/updateBookingStatus', async ({ id, status }, { rejectWithValue }) => {
  try {
    return await api.updateBookingStatus(id, status);
  } catch (error) {
    return handleError(error, rejectWithValue, 'Failed to update booking status');
  }
});

const marketplaceSlice = createSlice({
  name: 'marketplace',
  initialState,
  reducers: {
    setSelectedListing: (state, action: PayloadAction<ServiceListing | null>) => {
      state.selectedListing = action.payload;
    },
    setSelectedBooking: (state, action: PayloadAction<Booking | null>) => {
      state.selectedBooking = action.payload;
    },
    setFilters: (state, action: PayloadAction<Partial<SearchFilter>>) => {
      state.filters = { ...state.filters, ...action.payload };
    },
    clearFilters: (state) => {
      state.filters = initialState.filters;
      state.results = null;
      state.query = '';
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    // Search Services
    builder
      .addCase(searchServices.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(searchServices.fulfilled, (state, action) => {
        state.loading = false;
        state.results = action.payload;
        state.listings = action.payload.services;
      })
      .addCase(searchServices.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to search services';
      });

    // Fetch Service Listings
    builder
      .addCase(fetchServiceListings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchServiceListings.fulfilled, (state, action) => {
        state.loading = false;
        state.listings = action.payload.data;
      })
      .addCase(fetchServiceListings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to fetch service listings';
      });

    // Create Service Listing
    builder
      .addCase(createServiceListing.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createServiceListing.fulfilled, (state, action) => {
        state.loading = false;
        state.listings.push(action.payload.data);
      })
      .addCase(createServiceListing.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to create service listing';
      });

    // Update Service Listing
    builder
      .addCase(updateServiceListing.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateServiceListing.fulfilled, (state, action) => {
        state.loading = false;
        const index = state.listings.findIndex((listing) => listing.id === action.payload.data.id);
        if (index !== -1) {
          state.listings[index] = action.payload.data;
        }
      })
      .addCase(updateServiceListing.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to update service listing';
      });

    // Delete Service Listing
    builder
      .addCase(deleteServiceListing.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(deleteServiceListing.fulfilled, (state, action) => {
        state.loading = false;
        state.listings = state.listings.filter((listing) => listing.id !== action.meta.arg);
      })
      .addCase(deleteServiceListing.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to delete service listing';
      });

    // Fetch Bookings
    builder
      .addCase(fetchBookings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchBookings.fulfilled, (state, action) => {
        state.loading = false;
        state.bookings = action.payload.data;
      })
      .addCase(fetchBookings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to fetch bookings';
      });

    // Create Booking
    builder
      .addCase(createBooking.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createBooking.fulfilled, (state, action) => {
        state.loading = false;
        state.bookings.push(action.payload.data);
      })
      .addCase(createBooking.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to create booking';
      });

    // Update Booking Status
    builder
      .addCase(updateBookingStatus.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateBookingStatus.fulfilled, (state, action) => {
        state.loading = false;
        const index = state.bookings.findIndex((booking) => booking.id === action.payload.data.id);
        if (index !== -1) {
          state.bookings[index] = action.payload.data;
        }
      })
      .addCase(updateBookingStatus.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload || 'Failed to update booking status';
      });
  },
});

// Selectors
export const selectSearchResults = (state: { marketplace: MarketplaceState }) => state.marketplace.results;

export const {
  setSelectedListing,
  setSelectedBooking,
  setFilters,
  clearFilters,
  clearError,
} = marketplaceSlice.actions;

export default marketplaceSlice.reducer;
