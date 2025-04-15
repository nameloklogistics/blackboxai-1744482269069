import React from 'react';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import { renderWithProviders } from '../../../utils/test-utils';
import BookingManagement from '../index';
import { fetchBookings, updateBookingStatus } from '../../../store/slices/marketplaceSlice';
import { Booking, RootState, MarketplaceState, ApiResponse } from '../../../types';

// Mock the redux actions with proper types
jest.mock('../../../store/slices/marketplaceSlice', () => ({
  __esModule: true,
  fetchBookings: jest.fn(),
  updateBookingStatus: jest.fn()
}));

describe('BookingManagement', () => {
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
    trackingInfo: [
      {
        status: 'PICKED_UP',
        description: 'Cargo picked up from origin',
        location: 'Singapore Port',
        timestamp: '2023-01-01T00:00:00Z'
      },
      {
        status: 'IN_TRANSIT',
        description: 'Cargo in transit',
        location: 'South China Sea',
        timestamp: '2023-01-02T00:00:00Z'
      }
    ],
    createdAt: '2023-01-01T00:00:00Z',
    updatedAt: '2023-01-01T00:00:00Z',
  };

  const mockMarketplaceState: MarketplaceState = {
    listings: [],
    bookings: [mockBooking],
    selectedListing: null,
    selectedBooking: null,
    loading: false,
    error: null,
    filters: {
      page: 1,
      pageSize: 10,
      minRate: 0, // Default value instead of undefined
      maxRate: 0, // Default value instead of undefined
      query: undefined,
      serviceType: undefined,
      shipmentMode: undefined,
      origin: undefined,
      destination: undefined,
      minPrice: undefined,
      maxPrice: undefined,
      routeType: undefined,
      transitTime: undefined,
      providerId: undefined,
      isActive: undefined,
      validFrom: undefined,
      validUntil: undefined,
      sortBy: undefined,
      sortOrder: undefined
    },
  };

  const initialState: Partial<RootState> = {
    marketplace: mockMarketplaceState,
  };

  beforeEach(() => {
    jest.clearAllMocks();
    // Cast the mock function to any to avoid TypeScript errors with complex Redux types
    (fetchBookings as any).mockResolvedValue({
      type: 'marketplace/fetchBookings/fulfilled',
      payload: { success: true, data: [mockBooking] },
    });
  });

  it('renders loading state initially', () => {
    const loadingState: Partial<RootState> = {
      marketplace: { ...mockMarketplaceState, loading: true, bookings: [] },
    };

    renderWithProviders(<BookingManagement />, loadingState);
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });

  it('fetches bookings on mount', () => {
    renderWithProviders(<BookingManagement />, initialState);
    expect(fetchBookings).toHaveBeenCalled();
  });

  it('displays booking information correctly', () => {
    renderWithProviders(<BookingManagement />, initialState);
    
    expect(screen.getByText(`Booking #${mockBooking.id}`)).toBeInTheDocument();
    expect(screen.getByText(`Service ID: ${mockBooking.serviceId}`)).toBeInTheDocument();
    expect(screen.getByText(`Weight: ${mockBooking.cargoDetails.weight} kg`)).toBeInTheDocument();
    expect(screen.getByText(`Volume: ${mockBooking.cargoDetails.volume} m³`)).toBeInTheDocument();
    expect(screen.getByText(`Type: ${mockBooking.cargoDetails.type}`)).toBeInTheDocument();
    expect(screen.getByText(`Payment: ${mockBooking.paymentAmount} LMT`)).toBeInTheDocument();
  });

  it('displays correct status chips', () => {
    renderWithProviders(<BookingManagement />, initialState);
    
    expect(screen.getByText(mockBooking.status)).toHaveClass('MuiChip-colorWarning');
    expect(screen.getByText(mockBooking.paymentStatus)).toHaveClass('MuiChip-colorError');
  });

  it('opens tracking dialog when clicking view tracking button', async () => {
    renderWithProviders(<BookingManagement />, initialState);
    
    const viewTrackingButton = screen.getByText('View Tracking');
    fireEvent.click(viewTrackingButton);

    expect(screen.getByText(`Shipment Tracking - Booking #${mockBooking.id}`)).toBeInTheDocument();
    expect(screen.getByText('PICKED_UP')).toBeInTheDocument();
    expect(screen.getByText('IN_TRANSIT')).toBeInTheDocument();
  });

  it('closes tracking dialog when clicking close button', () => {
    renderWithProviders(<BookingManagement />, initialState);
    
    // Open dialog
    const viewTrackingButton = screen.getByText('View Tracking');
    fireEvent.click(viewTrackingButton);

    // Close dialog
    const closeButton = screen.getByText('Close');
    fireEvent.click(closeButton);

    expect(screen.queryByText(`Shipment Tracking - Booking #${mockBooking.id}`)).not.toBeInTheDocument();
  });

  it('handles error states gracefully', async () => {
    const errorMessage = 'Failed to fetch bookings';
    // Cast the mock function to any to avoid TypeScript errors with complex Redux types
    (fetchBookings as any).mockRejectedValue(new Error(errorMessage));
    
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    
    renderWithProviders(<BookingManagement />, {
      marketplace: { ...mockMarketplaceState, bookings: [] },
    });

    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalled();
    });

    consoleSpy.mockRestore();
  });

  it('displays empty state when no bookings exist', () => {
    renderWithProviders(<BookingManagement />, {
      marketplace: { ...mockMarketplaceState, bookings: [] },
    });

    expect(screen.queryByText(/Booking #/)).not.toBeInTheDocument();
  });
});
