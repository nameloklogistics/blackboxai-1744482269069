import Layout from './Layout';

export { Layout };
export default Layout;

// Export any additional layout-related components or types here
export interface LayoutProps {
  children: React.ReactNode;
}

// Export layout-specific constants
export const DRAWER_WIDTH = 240;
export const MOBILE_DRAWER_WIDTH = 280;
export const APP_BAR_HEIGHT = 64;

// Export layout-specific utilities
export const isNavLinkActive = (currentPath: string, linkPath: string): boolean => {
  if (linkPath === '/') {
    return currentPath === '/';
  }
  return currentPath.startsWith(linkPath);
};

// Export layout-specific types
export interface NavItem {
  label: string;
  path: string;
  icon?: React.ReactNode;
  children?: NavItem[];
  roles?: string[];
}

// Export layout-specific configurations
export const NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    path: '/dashboard',
  },
  {
    label: 'Marketplace',
    path: '/marketplace',
    children: [
      { label: 'Service Listings', path: '/marketplace/services' },
      { label: 'My Listings', path: '/marketplace/my-listings' },
    ],
  },
  {
    label: 'Shipments',
    path: '/shipments',
    children: [
      { label: 'Active Shipments', path: '/shipments/active' },
      { label: 'Completed Shipments', path: '/shipments/completed' },
    ],
  },
  {
    label: 'Token Management',
    path: '/tokens',
    children: [
      { label: 'Wallet', path: '/tokens/wallet' },
      { label: 'Transactions', path: '/tokens/transactions' },
      { label: 'Escrow', path: '/tokens/escrow' },
    ],
  },
  {
    label: 'Analytics',
    path: '/analytics',
    children: [
      { label: 'Overview', path: '/analytics/overview' },
      { label: 'Reports', path: '/analytics/reports' },
    ],
  },
];

// Export layout-specific hooks
export const useLayoutConfig = () => {
  return {
    drawerWidth: DRAWER_WIDTH,
    mobileDrawerWidth: MOBILE_DRAWER_WIDTH,
    appBarHeight: APP_BAR_HEIGHT,
    navItems: NAV_ITEMS,
  };
};
