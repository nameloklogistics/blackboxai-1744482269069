import { configureStore } from '@reduxjs/toolkit';
import marketplaceReducer, {
  fetchServiceListings,
  fetchBookings,
  createServiceListing,
  updateServiceListing,
  deleteServiceListing,
} from '../marketplaceSlice';
import { ApiResponse, ServiceListing, Booking } from '../../../types';
import { api } from '../../../services/api';

// Mock the API service
jest.mock('../../../services/api', () => ({
  api: {
    getServiceListings: jest.fn(),
    getBookings: jest.fn(),
    createServiceListing: jest.fn(),
    updateServiceListing: jest.fn(),
    deleteServiceListing: jest.fn(),
  },
}));

describe('marketplaceSlice', () => {
  const initialState = {
    listings: [],
    bookings: [],
    selectedListing: null,
    selectedBooking: null,
    loading: false,
    error: null,
    filters: {},
  };

  let store = configureStore({
    reducer: {
      marketplace: marketplaceReducer,
    },
    preloadedState: {
      marketplace: initialState,
    },
  });

  beforeEach(() => {
    store = configureStore({
      reducer: {
        marketplace: marketplaceReducer,
      },
      preloadedState: {
        marketplace: initialState,
      },
    });
    jest.clearAllMocks();
  });

  describe('initial state', () => {
    it('should handle initial state', () => {
      expect(store.getState().marketplace).toEqual(initialState);
    });
  });

  describe('fetchServiceListings', () => {
    const mockListings: ServiceListing[] = [
      { id: '1', providerId: '2', providerName: 'Test Provider', serviceType: 'FREIGHT_FORWARDING', shipmentMode: 'SEA', origin: 'Singapore', destination: 'Hong Kong', rate: 1000, description: 'Test service listing', isActive: true, createdAt: '2023-01-01T00:00:00Z', updatedAt: '2023-01-01T00:00:00Z' },
    ];

    it('should handle successful service listings fetch', async () => {
      (api.getServiceListings as jest.Mock).mockResolvedValue({
        success: true,
        data: mockListings,
      });

      const result = await store.dispatch(fetchServiceListings());
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/fetchServiceListings/fulfilled');
      expect(result.payload).toEqual(mockListings);
      expect(state.listings).toEqual(mockListings);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed service listings fetch', async () => {
      const errorMessage = 'Failed to fetch service listings';
      (api.getServiceListings as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchServiceListings());
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/fetchServiceListings/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });

  describe('fetchBookings', () => {
    const mockBookings: Booking[] = [
      { id: '1', serviceId: '1', customerId: '1', providerId: '2', status: 'PENDING', cargoDetails: { weight: 1000, volume: 10, type: 'General', description: 'Test cargo', hazardous: false }, paymentStatus: 'UNPAID', paymentAmount: 1000, trackingInfo: [], createdAt: '2023-01-01T00:00:00Z', updatedAt: '2023-01-01T00:00:00Z' },
    ];

    it('should handle successful bookings fetch', async () => {
      (api.getBookings as jest.Mock).mockResolvedValue({
        success: true,
        data: mockBookings,
      });

      const result = await store.dispatch(fetchBookings());
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/fetchBookings/fulfilled');
      expect(result.payload).toEqual(mockBookings);
      expect(state.bookings).toEqual(mockBookings);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed bookings fetch', async () => {
      const errorMessage = 'Failed to fetch bookings';
      (api.getBookings as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchBookings());
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/fetchBookings/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });

  describe('createServiceListing', () => {
    const mockNewListing: ServiceListing = {
      id: '2',
      providerId: '3',
      providerName: 'New Provider',
      serviceType: 'FREIGHT_FORWARDING',
      shipmentMode: 'AIR',
      origin: 'Tokyo',
      destination: 'Los Angeles',
      rate: 1500,
      description: 'New service listing',
      isActive: true,
      createdAt: '2023-01-01T00:00:00Z',
      updatedAt: '2023-01-01T00:00:00Z',
    };

    it('should handle successful service listing creation', async () => {
      (api.createServiceListing as jest.Mock).mockResolvedValue({
        success: true,
        data: mockNewListing,
      });

      const result = await store.dispatch(createServiceListing(mockNewListing));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/createServiceListing/fulfilled');
      expect(result.payload).toEqual(mockNewListing);
      expect(state.listings).toContainEqual(mockNewListing);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed service listing creation', async () => {
      const errorMessage = 'Failed to create service listing';
      (api.createServiceListing as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(createServiceListing(mockNewListing));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/createServiceListing/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });

  describe('updateServiceListing', () => {
    const mockUpdatedListing: ServiceListing = {
      id: '1',
      providerId: '2',
      providerName: 'Updated Provider',
      serviceType: 'FREIGHT_FORWARDING',
      shipmentMode: 'SEA',
      origin: 'Singapore',
      destination: 'Hong Kong',
      rate: 1200,
      description: 'Updated service listing',
      isActive: true,
      createdAt: '2023-01-01T00:00:00Z',
      updatedAt: '2023-01-01T00:00:00Z',
    };

    beforeEach(() => {
      store = configureStore({
        reducer: {
          marketplace: marketplaceReducer,
        },
        preloadedState: {
          marketplace: {
            ...initialState,
            listings: [{ ...mockUpdatedListing, providerName: 'Old Provider' }],
          },
        },
      });
    });

    it('should handle successful service listing update', async () => {
      (api.updateServiceListing as jest.Mock).mockResolvedValue({
        success: true,
        data: mockUpdatedListing,
      });

      const result = await store.dispatch(updateServiceListing({ id: '1', updates: mockUpdatedListing }));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/updateServiceListing/fulfilled');
      expect(result.payload).toEqual(mockUpdatedListing);
      expect(state.listings).toContainEqual(mockUpdatedListing);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed service listing update', async () => {
      const errorMessage = 'Failed to update service listing';
      (api.updateServiceListing as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(updateServiceListing({ id: '1', updates: mockUpdatedListing }));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/updateServiceListing/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });

  describe('deleteServiceListing', () => {
    const mockListingId = '1';
    const mockListing: ServiceListing = {
      id: mockListingId,
      providerId: '2',
      providerName: 'Test Provider',
      serviceType: 'FREIGHT_FORWARDING',
      shipmentMode: 'SEA',
      origin: 'Singapore',
      destination: 'Hong Kong',
      rate: 1000,
      description: 'Test service listing',
      isActive: true,
      createdAt: '2023-01-01T00:00:00Z',
      updatedAt: '2023-01-01T00:00:00Z',
    };

    beforeEach(() => {
      store = configureStore({
        reducer: {
          marketplace: marketplaceReducer,
        },
        preloadedState: {
          marketplace: {
            ...initialState,
            listings: [mockListing],
          },
        },
      });
    });

    it('should handle successful service listing deletion', async () => {
      (api.deleteServiceListing as jest.Mock).mockResolvedValue({ success: true });

      const result = await store.dispatch(deleteServiceListing(mockListingId));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/deleteServiceListing/fulfilled');
      expect(state.listings).not.toContainEqual(expect.objectContaining({ id: mockListingId }));
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed service listing deletion', async () => {
      const errorMessage = 'Failed to delete service listing';
      (api.deleteServiceListing as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(deleteServiceListing(mockListingId));
      const state = store.getState().marketplace;

      expect(result.type).toBe('marketplace/deleteServiceListing/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });
});
