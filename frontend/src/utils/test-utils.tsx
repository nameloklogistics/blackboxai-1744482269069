import React, { PropsWithChildren } from 'react';
import { render as rtlRender } from '@testing-library/react';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { ThemeProvider } from '@mui/material/styles';
import { BrowserRouter } from 'react-router-dom';
import { NotificationProvider } from '../components/shared/NotificationCenter';
import theme from '../theme';
import { rootReducer, RootState } from '../store';

// Mock data
export const mockUser = {
  id: '1',
  email: 'test@example.com',
  name: 'Test User',
  role: 'SHIPPER',
  walletAddress: 'GXXXXXXXXXXXXXXXXXXXXXX',
  membershipStatus: 'ACTIVE',
  createdAt: '2023-01-01T00:00:00Z',
  updatedAt: '2023-01-01T00:00:00Z',
};

export const mockServiceListing = {
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

export const mockBooking = {
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

export const mockTokenTransaction = {
  id: '1',
  type: 'PAYMENT',
  from: 'GXXXXXXXXXXXXXXXXXXXXXX',
  to: 'GXXXXXXXXXXXXXXXXXXXXXX',
  amount: '1000',
  status: 'COMPLETED',
  timestamp: '2023-01-01T00:00:00Z',
};

// Initial state type
interface InitialState {
  auth?: Partial<RootState['auth']>;
  marketplace?: Partial<RootState['marketplace']>;
  token?: Partial<RootState['token']>;
}

// Custom render function
function render(
  ui: React.ReactElement,
  {
    initialState = {},
    store = configureStore({
      reducer: rootReducer,
      preloadedState: initialState,
    }),
    ...renderOptions
  } = {}
) {
  function Wrapper({ children }: PropsWithChildren<{}>): JSX.Element {
    return (
      <Provider store={store}>
        <ThemeProvider theme={theme}>
          <BrowserRouter>
            <NotificationProvider>
              {children}
            </NotificationProvider>
          </BrowserRouter>
        </ThemeProvider>
      </Provider>
    );
  }
  return rtlRender(ui, { wrapper: Wrapper, ...renderOptions });
}

// Custom testing utilities
export function createMockStore(initialState: InitialState = {}) {
  return configureStore({
    reducer: rootReducer,
    preloadedState: initialState,
  });
}

export function renderWithProviders(
  ui: React.ReactElement,
  initialState: InitialState = {}
) {
  const store = createMockStore(initialState);
  return {
    ...render(ui, { store }),
    store,
  };
}

// Mock API response creator
export function createApiResponse<T>(data: T, success = true, message?: string) {
  return {
    success,
    data,
    message,
  };
}

// Mock error response creator
export function createApiError(
  status = 400,
  message = 'Bad Request',
  errors: string[] = []
) {
  return {
    response: {
      status,
      data: {
        success: false,
        message,
        errors,
      },
    },
  };
}

// Mock event handlers
export const mockHandlers = {
  click: jest.fn(),
  submit: jest.fn(),
  change: jest.fn(),
};

// Form helpers
export function fillForm(form: HTMLElement, data: Record<string, string>) {
  Object.entries(data).forEach(([name, value]) => {
    const input = form.querySelector(`[name="${name}"]`) as HTMLInputElement;
    if (input) {
      fireEvent.change(input, { target: { value } });
    }
  });
}

// Mock date for consistent testing
export function mockDate(isoDate: string) {
  const RealDate = Date;
  const mockDate = new RealDate(isoDate);
  
  global.Date = class extends RealDate {
    constructor(date?: string | number | Date) {
      super(date || mockDate);
    }
    
    static now() {
      return mockDate.getTime();
    }
  } as DateConstructor;
  
  return () => {
    global.Date = RealDate;
  };
}

// Re-export everything from RTL
export * from '@testing-library/react';
export { default as userEvent } from '@testing-library/user-event';
