import { configureStore } from '@reduxjs/toolkit';
import tokenReducer, {
  fetchTokenBalance,
  fetchTokenTransactions,
  transferTokens,
  createEscrow,
  releaseEscrow,
} from '../tokenSlice';
import { ApiResponse, TokenTransaction } from '../../../types';
import { api } from '../../../services/api';

// Mock the API service
jest.mock('../../../services/api', () => ({
  api: {
    getTokenBalance: jest.fn(),
    getTokenTransactions: jest.fn(),
    transferTokens: jest.fn(),
    createEscrow: jest.fn(),
    releaseEscrow: jest.fn(),
  },
}));

// Mock the blockchain service
jest.mock('../../../services/blockchain');

// Mock the constants
jest.mock('../../../constants', () => ({
  BLOCKCHAIN_NETWORK: 'TESTNET',
  BLOCKCHAIN_HORIZON_URL: 'https://horizon-testnet.stellar.org',
  TOKEN_CODE: 'TEST',
  TOKEN_ISSUER: 'GXXXXXXXXXXXXXXXXXXXXXX',
}));

describe('tokenSlice', () => {
  const initialState = {
    balance: '0',
    transactions: [],
    escrowBalances: [],
    loading: false,
    error: null,
  };

  let store = configureStore({
    reducer: {
      token: tokenReducer,
    },
    preloadedState: {
      token: initialState,
    },
  });

  beforeEach(() => {
    store = configureStore({
      reducer: {
        token: tokenReducer,
      },
      preloadedState: {
        token: initialState,
      },
    });
    jest.clearAllMocks();
  });

  describe('initial state', () => {
    it('should handle initial state', () => {
      expect(store.getState().token).toEqual(initialState);
    });
  });

  describe('fetchTokenBalance', () => {
    const mockWalletAddress = 'GXXXXXXXXXXXXXXXXXXXXXX';
    const mockBalance = '1000';

    it('should handle successful token balance fetch', async () => {
      (api.getTokenBalance as jest.Mock).mockResolvedValue({
        success: true,
        data: mockBalance,
      });

      const result = await store.dispatch(fetchTokenBalance(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/fulfilled');
      expect(result.payload).toBe(mockBalance);
      expect(state.balance).toBe(mockBalance);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
      expect(api.getTokenBalance).toHaveBeenCalledWith(mockWalletAddress);
    });

    it('should handle failed token balance fetch', async () => {
      const errorMessage = 'Network error. Please check your internet connection.';
      (api.getTokenBalance as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchTokenBalance(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
      expect(state.balance).toBe('0');
    });

    it('should set loading state while fetching', () => {
      (api.getTokenBalance as jest.Mock).mockImplementation(
        () => new Promise((resolve) => setTimeout(resolve, 100))
      );

      store.dispatch(fetchTokenBalance(mockWalletAddress));
      expect(store.getState().token.loading).toBe(true);
    });

    it('should handle API error response', async () => {
      const errorMessage = 'Network error. Please check your internet connection.';
      (api.getTokenBalance as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchTokenBalance(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });

  describe('edge cases', () => {
    const mockWalletAddress = 'GXXXXXXXXXXXXXXXXXXXXXX';

    it('should handle empty wallet address', async () => {
      const result = await store.dispatch(fetchTokenBalance(''));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/rejected');
      expect(state.error).toBeTruthy();
    });

    it('should handle network errors', async () => {
      const errorMessage = 'Network error. Please check your internet connection.';
      (api.getTokenBalance as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchTokenBalance(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/rejected');
      expect(state.error).toBe(errorMessage);
    });

    it('should handle malformed API responses', async () => {
      (api.getTokenBalance as jest.Mock).mockRejectedValue(new Error('Invalid response format'));

      const result = await store.dispatch(fetchTokenBalance(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchBalance/rejected');
      expect(state.error).toBeTruthy();
    });
  });

  describe('fetchTokenTransactions', () => {
    const mockWalletAddress = 'GXXXXXXXXXXXXXXXXXXXXXX';
    const mockTransactions: TokenTransaction[] = [
      {
        id: '1',
        from: 'GXXXXXXXXXXXXXXXXXXXXXX',
        to: 'GXXXXXXXXXXXXXXXXXXXXXY',
        amount: '100',
        status: 'COMPLETED',
        type: 'PAYMENT',
        timestamp: '2023-01-01T00:00:00Z',
      },
    ];

    it('should handle successful transactions fetch', async () => {
      (api.getTokenTransactions as jest.Mock).mockResolvedValue({
        success: true,
        data: mockTransactions,
      });

      const result = await store.dispatch(fetchTokenTransactions(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchTransactions/fulfilled');
      expect(result.payload).toEqual(mockTransactions);
      expect(state.transactions).toEqual(mockTransactions);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should handle failed transactions fetch', async () => {
      const errorMessage = 'Failed to fetch transactions';
      (api.getTokenTransactions as jest.Mock).mockRejectedValue(new Error(errorMessage));

      const result = await store.dispatch(fetchTokenTransactions(mockWalletAddress));
      const state = store.getState().token;

      expect(result.type).toBe('token/fetchTransactions/rejected');
      expect(state.loading).toBe(false);
      expect(state.error).toBe(errorMessage);
    });
  });
});
