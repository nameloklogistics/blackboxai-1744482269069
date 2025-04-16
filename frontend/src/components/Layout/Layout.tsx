import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../../types';
import {
  Box,
  Drawer,
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  Avatar,
  Menu,
  MenuItem,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard,
  Store,
  LocalShipping,
  AccountBalanceWallet,
  Person,
  Settings,
  ExitToApp,
  ChevronLeft,
} from '@mui/icons-material';
import { useNotification } from '../shared/NotificationCenter';
import { ROUTES, USER_ROLES } from '../../constants';

const drawerWidth = 240;

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const { showSuccess } = useNotification();

  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const { user } = useSelector((state: RootState) => state.auth);
  const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    // Implement logout logic here
    handleMenuClose();
    showSuccess('Successfully logged out');
    navigate(ROUTES.LOGIN);
  };

  const navigationItems = [
    {
      text: 'Dashboard',
      icon: <Dashboard />,
      path: ROUTES.DASHBOARD,
      roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER, USER_ROLES.SHIPPER, USER_ROLES.ADMIN],
    },
    {
      text: 'Marketplace',
      icon: <Store />,
      path: ROUTES.MARKETPLACE.ROOT,
      roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER, USER_ROLES.SHIPPER],
    },
    {
      text: 'Bookings',
      icon: <LocalShipping />,
      path: ROUTES.BOOKINGS.ROOT,
      roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER, USER_ROLES.SHIPPER],
    },
    {
      text: 'Token Management',
      icon: <AccountBalanceWallet />,
      path: ROUTES.TOKENS.ROOT,
      roles: [USER_ROLES.FREIGHT_FORWARDER, USER_ROLES.CUSTOMS_BROKER, USER_ROLES.SHIPPER],
    },
  ];

  const drawer = (
    <Box>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          p: 2,
        }}
      >
        <Typography variant="h6" noWrap component="div">
          Logistics Platform
        </Typography>
      </Box>
      <Divider />
      <List>
        {navigationItems.map(
          (item) =>
            user &&
            item.roles.includes(user.role) && (
              <ListItem
                button
                key={item.text}
                onClick={() => {
                  navigate(item.path);
                  if (isMobile) handleDrawerToggle();
                }}
                selected={location.pathname.startsWith(item.path)}
                sx={{
                  '&.Mui-selected': {
                    backgroundColor: theme.palette.primary.light,
                    '&:hover': {
                      backgroundColor: theme.palette.primary.light,
                    },
                  },
                }}
              >
                <ListItemIcon sx={{ color: location.pathname.startsWith(item.path) ? 'primary.main' : 'inherit' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText primary={item.text} />
              </ListItem>
            )
        )}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Box sx={{ flexGrow: 1 }} />
          <IconButton
            onClick={handleMenuOpen}
            size="small"
            sx={{ ml: 2 }}
            aria-controls="user-menu"
            aria-haspopup="true"
          >
            <Avatar sx={{ width: 32, height: 32 }}>
              {user?.name?.charAt(0) || 'U'}
            </Avatar>
          </IconButton>
          <Menu
            id="user-menu"
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            onClick={handleMenuClose}
          >
            <MenuItem onClick={() => navigate(ROUTES.PROFILE.ROOT)}>
              <ListItemIcon>
                <Person fontSize="small" />
              </ListItemIcon>
              Profile
            </MenuItem>
            <MenuItem onClick={() => navigate(ROUTES.PROFILE.SETTINGS)}>
              <ListItemIcon>
                <Settings fontSize="small" />
              </ListItemIcon>
              Settings
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <ListItemIcon>
                <ExitToApp fontSize="small" />
              </ListItemIcon>
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
      >
        {isMobile ? (
          <Drawer
            variant="temporary"
            anchor="left"
            open={mobileOpen}
            onClose={handleDrawerToggle}
            ModalProps={{
              keepMounted: true, // Better open performance on mobile.
            }}
            sx={{
              display: { xs: 'block', md: 'none' },
              '& .MuiDrawer-paper': {
                boxSizing: 'border-box',
                width: drawerWidth,
              },
            }}
          >
            {drawer}
          </Drawer>
        ) : (
          <Drawer
            variant="permanent"
            sx={{
              display: { xs: 'none', md: 'block' },
              '& .MuiDrawer-paper': {
                boxSizing: 'border-box',
                width: drawerWidth,
              },
            }}
            open
          >
            {drawer}
          </Drawer>
        )}
      </Box>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${drawerWidth}px)` },
          mt: '64px',
        }}
      >
        {children}
      </Box>
    </Box>
  );
};

export default Layout;
