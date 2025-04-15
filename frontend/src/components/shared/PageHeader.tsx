import React from 'react';
import {
  Box,
  Typography,
  Breadcrumbs,
  Link,
  Stack,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { Link as RouterLink, useLocation } from 'react-router-dom';
import { ROUTES } from '../../constants';

export interface PageHeaderProps {
  title: string;
  subtitle?: string;
  breadcrumbs?: Array<{
    label: string;
    path?: string;
  }>;
  actions?: React.ReactNode;
}

const PageHeader: React.FC<PageHeaderProps> = ({
  title,
  subtitle,
  breadcrumbs,
  actions,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const location = useLocation();

  // Generate default breadcrumbs based on current path if not provided
  const defaultBreadcrumbs = React.useMemo(() => {
    if (breadcrumbs) return breadcrumbs;

    const paths = location.pathname.split('/').filter(Boolean);
    return paths.map((path, index) => {
      const label = path.charAt(0).toUpperCase() + path.slice(1).toLowerCase();
      const fullPath = `/${paths.slice(0, index + 1).join('/')}`;
      return {
        label,
        path: fullPath,
      };
    });
  }, [location.pathname, breadcrumbs]);

  return (
    <Box
      sx={{
        mb: 3,
        px: { xs: 2, sm: 3 },
        py: 2,
        backgroundColor: 'background.paper',
        borderRadius: 1,
        boxShadow: 1,
      }}
    >
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        justifyContent="space-between"
        alignItems={{ xs: 'flex-start', sm: 'center' }}
        spacing={2}
      >
        <Box>
          {defaultBreadcrumbs.length > 0 && (
            <Breadcrumbs
              aria-label="breadcrumb"
              sx={{ mb: 1, '& .MuiBreadcrumbs-separator': { mx: 1 } }}
            >
              <Link
                component={RouterLink}
                to={ROUTES.DASHBOARD}
                color="inherit"
                sx={{
                  textDecoration: 'none',
                  '&:hover': { textDecoration: 'underline' },
                }}
              >
                Home
              </Link>
              {defaultBreadcrumbs.map((crumb, index) => {
                const isLast = index === defaultBreadcrumbs.length - 1;
                return isLast ? (
                  <Typography
                    key={crumb.label}
                    color="text.primary"
                    sx={{ fontWeight: 'medium' }}
                  >
                    {crumb.label}
                  </Typography>
                ) : (
                  <Link
                    key={crumb.label}
                    component={RouterLink}
                    to={crumb.path || '#'}
                    color="inherit"
                    sx={{
                      textDecoration: 'none',
                      '&:hover': { textDecoration: 'underline' },
                    }}
                  >
                    {crumb.label}
                  </Link>
                );
              })}
            </Breadcrumbs>
          )}
          <Typography
            variant={isMobile ? 'h6' : 'h5'}
            component="h1"
            sx={{
              fontWeight: 'bold',
              color: 'text.primary',
              mb: subtitle ? 0.5 : 0,
            }}
          >
            {title}
          </Typography>
          {subtitle && (
            <Typography
              variant="body1"
              color="text.secondary"
              sx={{ mt: 0.5 }}
            >
              {subtitle}
            </Typography>
          )}
        </Box>
        {actions && (
          <Box
            sx={{
              display: 'flex',
              gap: 1,
              flexWrap: 'wrap',
              justifyContent: { xs: 'flex-start', sm: 'flex-end' },
              width: { xs: '100%', sm: 'auto' },
            }}
          >
            {actions}
          </Box>
        )}
      </Stack>
    </Box>
  );
};

export default PageHeader;

// Example usage:
/*
import PageHeader from '@components/shared/PageHeader';
import { Button } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';

const MyPage = () => {
  return (
    <>
      <PageHeader
        title="Service Listings"
        subtitle="Manage your service listings"
        breadcrumbs={[
          { label: 'Marketplace', path: '/marketplace' },
          { label: 'Service Listings' },
        ]}
        actions={
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {}}
          >
            Add New Listing
          </Button>
        }
      />
      {/* Page content */}
    </>
  );
};
*/
