import { configureStore, ThunkAction, Action, ThunkDispatch, AnyAction } from '@reduxjs/toolkit';
import { useDispatch, useSelector, TypedUseSelectorHook } from 'react-redux';
import authReducer from './slices/authSlice';
import marketplaceReducer from './slices/marketplaceSlice';
import tokenReducer from './slices/tokenSlice';

// Configure middleware
const middleware = (getDefaultMiddleware: any) =>
  getDefaultMiddleware({
    serializableCheck: {
      // Ignore these action types
      ignoredActions: ['auth/login/fulfilled', 'auth/register/fulfilled'],
      // Ignore these field paths in all actions
      ignoredActionPaths: ['payload.token'],
      // Ignore these paths in the state
      ignoredPaths: ['auth.token'],
    },
  });

// Create store
export const store = configureStore({
  reducer: {
    auth: authReducer,
    marketplace: marketplaceReducer,
    token: tokenReducer,
  },
  middleware,
});

// Export types
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = ThunkDispatch<RootState, unknown, AnyAction>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;

// Create typed hooks
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

// Create API middleware
export const api = {
  setAuthToken: (token: string | null) => {
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  },
};

// Initialize auth token from localStorage
const token = localStorage.getItem('token');
if (token) {
  api.setAuthToken(token);
}

export default store;
