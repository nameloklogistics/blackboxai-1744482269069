import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../../utils/test-utils';
import Dashboard from '../index';
import { fetchTokenBalance } from '../../../store/slices/tokenSlice';
import { fetchServiceListings, fetchBookings } from '../../../store/slices/marketplaceSlice';
import { 
  User, 
  ServiceListing, 
  Booking, 
  ApiResponse,
  RootState,
  AuthState,
  TokenState,
  MarketplaceState
} from '../../../types';

// Mock the redux actions
jest.mock('../../../store/slices/tokenSlice');
jest.mock('../../../store/slices/marketplaceSlice');

describe('Dashboard', () => {
  // Type-safe mock data
  const mockUser: User = {
    id: '1',
    email: 'test@example.com',
    name: 'Test User',
    role: 'SHIPPER',
    walletAddress: 'GXXXXXXXXXXXXXXXXXXXXXX',
    membershipStatus: 'ACTIVE',
    createdAt: '2023-01-01T00:00:00Z',
    updatedAt: '2023-01-01T00:00:00Z',
  };

  const mockServiceListing: ServiceListing = {
    id: '1',
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

  const mockBooking: Booking = {
    id: '1',
    serviceId: '1',
    customerId: '1',
    providerId: '2',
    status: 'PENDING',
    cargoDetails: {
      weight: 1000,
      volume: 10,
      type: 'General',
      description: 'Test cargo',
      hazardous: false,
    },
    paymentStatus: 'UNPAID',
    paymentAmount: 1000,
    trackingInfo: [],
    createdAt: '2023-01-01T00:00:00Z',
    updatedAt: '2023-01-01T00:00:00Z',
  };

  const mockActiveBooking: Booking = {
    ...mockBooking,
    status: 'CONFIRMED',
  };

  const mockCompletedBooking: Booking = {
    ...mockBooking,
    status: 'COMPLETED',
  };

  const mockAuthState: AuthState = {
    user: mockUser,
    token: null,
    isAuthenticated: true,
    loading: false,
    error: null,
  };

  const mockTokenState: TokenState = {
    balance: '1000',
    transactions: [],
    escrowBalances: [],
    loading: false,
    error: null,
  };

  const mockMarketplaceState: MarketplaceState = {
    listings: [mockServiceListing],
    bookings: [mockActiveBooking, mockCompletedBooking],
    selectedListing: null,
    selectedBooking: null,
    loading: false,
    error: null,
    filters: {},
  };

  const initialState: Partial<RootState> = {
    auth: mockAuthState,
    token: mockTokenState,
    marketplace: mockMarketplaceState,
  };

  beforeEach(() => {
    jest.clearAllMocks();
    (fetchTokenBalance as unknown as jest.Mock).mockResolvedValue({ 
      type: 'token/fetchBalance/fulfilled', 
      payload: '1000' 
    });
    (fetchServiceListings as unknown as jest.Mock).mockResolvedValue({ 
      type: 'marketplace/fetchListings/fulfilled', 
      payload: { success: true, data: [mockServiceListing] } as ApiResponse<ServiceListing[]>
    });
    (fetchBookings as unknown as jest.Mock).mockResolvedValue({ 
      type: 'marketplace/fetchBookings/fulfilled', 
      payload: { success: true, data: [mockActiveBooking, mockCompletedBooking] } as ApiResponse<Booking[]>
    });
  });

  it('renders loading state initially', () => {
    const loadingState: Partial<RootState> = {
      auth: mockAuthState,
      token: { ...mockTokenState, loading: true },
      marketplace: { ...mockMarketplaceState, loading: true, listings: [], bookings: [] },
    };

    renderWithProviders(<Dashboard />, loadingState);
    expect(screen.getAllByRole('progressbar')).toHaveLength(4);
  });

  it('fetches data on mount', () => {
    const baseState: Partial<RootState> = {
      auth: mockAuthState,
    };

    renderWithProviders(<Dashboard />, baseState);

    expect(fetchTokenBalance).toHaveBeenCalledWith(mockUser.walletAddress);
    expect(fetchServiceListings).toHaveBeenCalled();
    expect(fetchBookings).toHaveBeenCalled();
  });

  it('displays user welcome message', () => {
    renderWithProviders(<Dashboard />, initialState);
    expect(screen.getByText(`Welcome back, ${mockUser.name}`)).toBeInTheDocument();
  });

  it('displays correct token balance', () => {
    renderWithProviders(<Dashboard />, initialState);
    expect(screen.getByText('1000')).toBeInTheDocument();
  });

  it('displays correct number of active bookings', () => {
    renderWithProviders(<Dashboard />, initialState);
    expect(screen.getByText('1')).toBeInTheDocument();
  });

  it('displays recent listings correctly', () => {
    renderWithProviders(<Dashboard />, initialState);
    expect(screen.getByText(mockServiceListing.providerName)).toBeInTheDocument();
    expect(screen.getByText(`${mockServiceListing.rate} LMT`)).toBeInTheDocument();
    expect(screen.getByText(`${mockServiceListing.origin} → ${mockServiceListing.destination}`)).toBeInTheDocument();
  });

  it('displays active bookings correctly', () => {
    renderWithProviders(<Dashboard />, initialState);
    expect(screen.getByText(`Booking #${mockActiveBooking.id}`)).toBeInTheDocument();
    expect(screen.getByText(`Status: ${mockActiveBooking.status}`)).toBeInTheDocument();
    expect(screen.getByText(`Payment: ${mockActiveBooking.paymentStatus}`)).toBeInTheDocument();
  });

  it('handles error states gracefully', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    
    (fetchTokenBalance as unknown as jest.Mock).mockRejectedValue(new Error('Failed to fetch balance'));
    (fetchServiceListings as unknown as jest.Mock).mockRejectedValue(new Error('Failed to fetch listings'));
    (fetchBookings as unknown as jest.Mock).mockRejectedValue(new Error('Failed to fetch bookings'));

    renderWithProviders(<Dashboard />, { auth: mockAuthState });

    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalled();
    });

    consoleSpy.mockRestore();
  });

  it('does not fetch token balance without wallet address', () => {
    // Create a new user object without walletAddress
    const userWithoutWallet = { ...mockUser } as User;
    delete userWithoutWallet.walletAddress;

    const stateWithoutWallet: Partial<RootState> = {
      auth: { 
        ...mockAuthState, 
        user: userWithoutWallet,
      },
    };

    renderWithProviders(<Dashboard />, stateWithoutWallet);

    expect(fetchTokenBalance).not.toHaveBeenCalled();
    expect(fetchServiceListings).toHaveBeenCalled();
    expect(fetchBookings).toHaveBeenCalled();
  });

  it('updates when new data is received', async () => {
    const { rerender } = renderWithProviders(<Dashboard />, initialState);

    const updatedState: Partial<RootState> = {
      ...initialState,
      token: {
        ...mockTokenState,
        balance: '2000',
      },
    };

    rerender(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('2000')).toBeInTheDocument();
    });
  });

  it('filters active bookings correctly', () => {
    const stateWithMultipleBookings: Partial<RootState> = {
      ...initialState,
      marketplace: {
        ...mockMarketplaceState,
        bookings: [
          { ...mockBooking, status: 'CONFIRMED', id: '1' },
          { ...mockBooking, status: 'IN_PROGRESS', id: '2' },
          { ...mockBooking, status: 'COMPLETED', id: '3' },
          { ...mockBooking, status: 'CANCELLED', id: '4' },
        ],
      },
    };

    renderWithProviders(<Dashboard />, stateWithMultipleBookings);

    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('Booking #1')).toBeInTheDocument();
    expect(screen.getByText('Booking #2')).toBeInTheDocument();
    expect(screen.queryByText('Booking #3')).not.toBeInTheDocument();
    expect(screen.queryByText('Booking #4')).not.toBeInTheDocument();
  });
});
