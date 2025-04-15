import React from 'react';
import { Navigate, RouteObject } from 'react-router-dom';
import { USER_ROLES, ROUTES } from '../constants';
import Layout from '../components/Layout';
import ProtectedRoute from './ProtectedRoute';

// Pages
const Dashboard = React.lazy(() => import('../pages/Dashboard'));
const ServiceListings = React.lazy(() => import('../pages/ServiceListings'));
const TokenManagement = React.lazy(() => import('../pages/TokenManagement'));
const BookingManagement = React.lazy(() => import('../pages/BookingManagement'));
const Login = React.lazy(() => import('../pages/Auth/Login'));
const Register = React.lazy(() => import('../pages/Auth/Register'));

// Route configurations with role-based access control
interface RouteConfig extends RouteObject {
  roles?: string[];
  children?: RouteConfig[];
}

export const routes: RouteConfig[] = [
  {
    path: ROUTES.LOGIN,
    element: <Login />,
  },
  {
    path: ROUTES.REGISTER,
    element: <Register />,
  },
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        path: '/',
        element: <Layout />,
        children: [
          {
            path: '/',
            element: <Navigate to={ROUTES.DASHBOARD} replace />,
          },
          {
            path: ROUTES.DASHBOARD,
            element: <Dashboard />,
          },
          {
            path: ROUTES.MARKETPLACE.ROOT,
            children: [
              {
                path: '',
                element: <ServiceListings />,
              },
              {
                path: 'my-listings',
                element: <ServiceListings />,
                roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER],
              },
              {
                path: 'create',
                element: <ServiceListings />,
                roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER],
              },
              {
                path: 'edit/:id',
                element: <ServiceListings />,
                roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER],
              },
            ],
          },
          {
            path: ROUTES.BOOKINGS.ROOT,
            children: [
              {
                path: '',
                element: <BookingManagement />,
              },
              {
                path: 'active',
                element: <BookingManagement />,
              },
              {
                path: 'completed',
                element: <BookingManagement />,
              },
              {
                path: ':id',
                element: <BookingManagement />,
              },
            ],
          },
          {
            path: ROUTES.TOKENS.ROOT,
            children: [
              {
                path: '',
                element: <TokenManagement />,
              },
              {
                path: 'wallet',
                element: <TokenManagement />,
              },
              {
                path: 'transactions',
                element: <TokenManagement />,
              },
              {
                path: 'escrow',
                element: <TokenManagement />,
              },
            ],
          },
          {
            path: ROUTES.PROFILE.ROOT,
            children: [
              {
                path: '',
                element: <Navigate to={ROUTES.PROFILE.SETTINGS} replace />,
              },
              {
                path: 'settings',
                element: <React.Fragment>Profile Settings</React.Fragment>,
              },
              {
                path: 'security',
                element: <React.Fragment>Security Settings</React.Fragment>,
              },
            ],
          },
          {
            path: ROUTES.ANALYTICS.ROOT,
            children: [
              {
                path: '',
                element: <Navigate to={ROUTES.ANALYTICS.OVERVIEW} replace />,
              },
              {
                path: 'overview',
                element: <React.Fragment>Analytics Overview</React.Fragment>,
                roles: [USER_ROLES.ADMIN],
              },
              {
                path: 'reports',
                element: <React.Fragment>Analytics Reports</React.Fragment>,
                roles: [USER_ROLES.ADMIN],
              },
            ],
          },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to={ROUTES.DASHBOARD} replace />,
  },
];

// Helper function to check if user has required role
export const hasRequiredRole = (userRole: string, requiredRoles?: string[]): boolean => {
  if (!requiredRoles || requiredRoles.length === 0) {
    return true;
  }
  return requiredRoles.includes(userRole);
};

// Helper function to get route title
export const getRouteTitle = (pathname: string): string => {
  const path = pathname.split('/').filter(Boolean);
  if (path.length === 0) return 'Dashboard';
  
  const title = path[path.length - 1]
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
  
  return title;
};

// Helper function to get breadcrumbs
export const getBreadcrumbs = (pathname: string) => {
  const paths = pathname.split('/').filter(Boolean);
  return paths.map((path, index) => ({
    label: path.charAt(0).toUpperCase() + path.slice(1),
    href: '/' + paths.slice(0, index + 1).join('/'),
  }));
};

export default routes;
