import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { TokenTransaction, EscrowBalance } from '../../types';
import { api } from '../../services/api';
import { blockchain } from '../../services/blockchain';

interface TokenState {
  balance: string;
  transactions: TokenTransaction[];
  escrowBalances: EscrowBalance[];
  loading: boolean;
  error: string | null;
}

const initialState: TokenState = {
  balance: '0',
  transactions: [],
  escrowBalances: [],
  loading: false,
  error: null,
};

// Async thunks
export const fetchTokenBalance = createAsyncThunk(
  'token/fetchBalance',
  async (address: string) => {
    const response = await api.getTokenBalance(address);
    return response.data;
  }
);

export const fetchTokenTransactions = createAsyncThunk(
  'token/fetchTransactions',
  async (address: string) => {
    const response = await api.getTokenTransactions(address);
    return response.data;
  }
);

export const transferTokens = createAsyncThunk(
  'token/transfer',
  async ({
    to,
    amount,
    memo,
  }: {
    to: string;
    amount: string;
    memo?: string;
  }) => {
    const response = await api.transferTokens(to, amount, memo);
    return response.data;
  }
);

export const createEscrow = createAsyncThunk(
  'token/createEscrow',
  async ({
    amount,
    bookingId,
    duration,
  }: {
    amount: string;
    bookingId: string;
    duration: number;
  }) => {
    const response = await api.createEscrow(amount, bookingId, duration);
    return response.data;
  }
);

export const releaseEscrow = createAsyncThunk(
  'token/releaseEscrow',
  async (escrowId: string) => {
    const response = await api.releaseEscrow(escrowId);
    return response.data;
  }
);

const tokenSlice = createSlice({
  name: 'token',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    updateBalance: (state, action) => {
      state.balance = action.payload;
    },
  },
  extraReducers: (builder) => {
    // Fetch Balance
    builder
      .addCase(fetchTokenBalance.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTokenBalance.fulfilled, (state, action) => {
        state.loading = false;
        state.balance = action.payload;
      })
      .addCase(fetchTokenBalance.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch balance';
      });

    // Fetch Transactions
    builder
      .addCase(fetchTokenTransactions.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTokenTransactions.fulfilled, (state, action) => {
        state.loading = false;
        state.transactions = action.payload;
      })
      .addCase(fetchTokenTransactions.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch transactions';
      });

    // Transfer Tokens
    builder
      .addCase(transferTokens.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(transferTokens.fulfilled, (state, action) => {
        state.loading = false;
        state.transactions = [action.payload, ...state.transactions];
        // Update balance after successful transfer
        blockchain.getBalance(action.payload.from).then((balance) => {
          state.balance = balance;
        });
      })
      .addCase(transferTokens.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to transfer tokens';
      });

    // Create Escrow
    builder
      .addCase(createEscrow.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createEscrow.fulfilled, (state, action) => {
        state.loading = false;
        state.escrowBalances = [...state.escrowBalances, action.payload];
      })
      .addCase(createEscrow.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to create escrow';
      });

    // Release Escrow
    builder
      .addCase(releaseEscrow.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(releaseEscrow.fulfilled, (state, action) => {
        state.loading = false;
        state.escrowBalances = state.escrowBalances.filter(
          (escrow) => escrow.id !== action.meta.arg
        );
        // Update balance after successful release
        blockchain.getBalance(action.payload.to).then((balance) => {
          state.balance = balance;
        });
      })
      .addCase(releaseEscrow.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to release escrow';
      });
  },
});

export const { clearError, updateBalance } = tokenSlice.actions;
export default tokenSlice.reducer;
