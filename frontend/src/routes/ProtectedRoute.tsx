import React, { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { CircularProgress, Box } from '@mui/material';
import { RootState } from '../store';
import { getCurrentUser } from '../store/slices/authSlice';
import { ROUTES } from '../constants';

const ProtectedRoute: React.FC = () => {
  const location = useLocation();
  const dispatch = useDispatch();
  const { isAuthenticated, loading, user } = useSelector((state: RootState) => state.auth);

  useEffect(() => {
    if (!user && isAuthenticated) {
      dispatch(getCurrentUser());
    }
  }, [dispatch, user, isAuthenticated]);

  // Show loading spinner while checking authentication
  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />;
  }

  // Check if user has required role for the route
  const route = location.pathname.split('/')[1];
  const requiredRoles = getRequiredRoles(route);

  if (requiredRoles && user && !requiredRoles.includes(user.role)) {
    return <Navigate to={ROUTES.DASHBOARD} replace />;
  }

  // Render child routes
  return <Outlet />;
};

// Helper function to determine required roles for each route
const getRequiredRoles = (route: string): string[] | null => {
  switch (route) {
    case 'marketplace':
      return ['FREIGHT_FORWARDER', 'CUSTOMS_BROKER', 'SHIPPER', 'ADMIN'];
    case 'bookings':
      return ['FREIGHT_FORWARDER', 'CUSTOMS_BROKER', 'SHIPPER', 'ADMIN'];
    case 'tokens':
      return ['FREIGHT_FORWARDER', 'CUSTOMS_BROKER', 'SHIPPER', 'ADMIN'];
    case 'analytics':
      return ['ADMIN'];
    default:
      return null;
  }
};

export default ProtectedRoute;
