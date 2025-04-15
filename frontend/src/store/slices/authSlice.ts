import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { User } from '../../types';
import { api } from '../../services/api';
import { storage } from '../../utils';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  user: null,
  token: storage.get('token'),
  isAuthenticated: Boolean(storage.get('token')),
  loading: false,
  error: null,
};

// Async thunks
export const login = createAsyncThunk(
  'auth/login',
  async ({ email, password }: { email: string; password: string }) => {
    const response = await api.login(email, password);
    storage.set('token', response.data.token);
    return response.data;
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async (userData: Partial<User>) => {
    const response = await api.register(userData);
    storage.set('token', response.data.token);
    return response.data;
  }
);

export const getCurrentUser = createAsyncThunk('auth/getCurrentUser', async () => {
  const response = await api.getCurrentUser();
  return response.data;
});

export const connectWallet = createAsyncThunk(
  'auth/connectWallet',
  async (secretKey: string) => {
    const response = await api.connectWallet(secretKey, 'signature');
    return response.data;
  }
);

export const disconnectWallet = createAsyncThunk(
  'auth/disconnectWallet',
  async () => {
    const response = await api.disconnectWallet();
    return response.data;
  }
);

export const updateProfile = createAsyncThunk(
  'auth/updateProfile',
  async (updates: Partial<User>) => {
    const response = await api.updateProfile(updates);
    return response.data;
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.error = null;
      storage.remove('token');
    },
    clearError: (state) => {
      state.error = null;
    },
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
      state.isAuthenticated = true;
    },
  },
  extraReducers: (builder) => {
    // Login
    builder
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.isAuthenticated = true;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Login failed';
      });

    // Register
    builder
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.isAuthenticated = true;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Registration failed';
      });

    // Get Current User
    builder
      .addCase(getCurrentUser.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getCurrentUser.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        state.isAuthenticated = true;
      })
      .addCase(getCurrentUser.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch user';
        state.isAuthenticated = false;
        storage.remove('token');
      });

    // Connect Wallet
    builder
      .addCase(connectWallet.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(connectWallet.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(connectWallet.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to connect wallet';
      });

    // Disconnect Wallet
    builder
      .addCase(disconnectWallet.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(disconnectWallet.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(disconnectWallet.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to disconnect wallet';
      });

    // Update Profile
    builder
      .addCase(updateProfile.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateProfile.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(updateProfile.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to update profile';
      });
  },
});

export const { logout, clearError, setUser } = authSlice.actions;
export default authSlice.reducer;
