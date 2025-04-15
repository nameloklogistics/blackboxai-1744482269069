import React, { useEffect } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  Button,
  CircularProgress,
} from '@mui/material';
import {
  Timeline,
  AccountBalanceWallet,
  LocalShipping,
  Assessment,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../store';
import { fetchTokenBalance } from '../../store/slices/tokenSlice';
import { fetchServiceListings, fetchBookings } from '../../store/slices/marketplaceSlice';

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const { balance, loading: tokenLoading } = useAppSelector((state) => state.token);
  const { listings, bookings, loading: marketplaceLoading } = useAppSelector(
    (state) => state.marketplace
  );

  useEffect(() => {
    if (user?.walletAddress) {
      dispatch(fetchTokenBalance(user.walletAddress));
    }
    dispatch(fetchServiceListings());
    dispatch(fetchBookings());
  }, [dispatch, user]);

  const activeBookings = bookings.filter(
    (booking) => booking.status === 'CONFIRMED' || booking.status === 'IN_PROGRESS'
  );

  const recentListings = listings.slice(0, 5);

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Grid container spacing={3}>
        {/* Welcome Section */}
        <Grid item xs={12}>
          <Typography variant="h4" gutterBottom>
            Welcome back, {user?.name}
          </Typography>
          <Typography variant="subtitle1" color="textSecondary">
            Here's what's happening with your logistics operations
          </Typography>
        </Grid>

        {/* Stats Cards */}
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <AccountBalanceWallet color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    LMT Balance
                  </Typography>
                  <Typography variant="h6">
                    {tokenLoading ? <CircularProgress size={20} /> : balance}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <LocalShipping color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Active Bookings
                  </Typography>
                  <Typography variant="h6">
                    {marketplaceLoading ? (
                      <CircularProgress size={20} />
                    ) : (
                      activeBookings.length
                    )}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <Timeline color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Total Transactions
                  </Typography>
                  <Typography variant="h6">
                    {marketplaceLoading ? (
                      <CircularProgress size={20} />
                    ) : (
                      bookings.length
                    )}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <Assessment color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Active Listings
                  </Typography>
                  <Typography variant="h6">
                    {marketplaceLoading ? (
                      <CircularProgress size={20} />
                    ) : (
                      listings.filter((l) => l.isActive).length
                    )}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Listings */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
              <Typography variant="h6">Recent Listings</Typography>
              <Button color="primary" variant="outlined" size="small">
                View All
              </Button>
            </Box>
            {marketplaceLoading ? (
              <Box display="flex" justifyContent="center" p={3}>
                <CircularProgress />
              </Box>
            ) : (
              recentListings.map((listing) => (
                <Card key={listing.id} sx={{ mb: 1 }}>
                  <CardContent>
                    <Typography variant="subtitle1">{listing.providerName}</Typography>
                    <Typography variant="body2" color="textSecondary">
                      {listing.serviceType} - {listing.shipmentMode}
                    </Typography>
                    <Typography variant="body2">
                      {listing.origin} → {listing.destination}
                    </Typography>
                    <Typography variant="h6" color="primary">
                      {listing.rate} LMT
                    </Typography>
                  </CardContent>
                </Card>
              ))
            )}
          </Paper>
        </Grid>

        {/* Active Bookings */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 2 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
              <Typography variant="h6">Active Bookings</Typography>
              <Button color="primary" variant="outlined" size="small">
                View All
              </Button>
            </Box>
            {marketplaceLoading ? (
              <Box display="flex" justifyContent="center" p={3}>
                <CircularProgress />
              </Box>
            ) : (
              activeBookings.map((booking) => (
                <Card key={booking.id} sx={{ mb: 1 }}>
                  <CardContent>
                    <Typography variant="subtitle1">Booking #{booking.id}</Typography>
                    <Typography variant="body2" color="textSecondary">
                      Status: {booking.status}
                    </Typography>
                    <Typography variant="body2">
                      Payment: {booking.paymentStatus}
                    </Typography>
                    <Typography variant="h6" color="primary">
                      {booking.paymentAmount} LMT
                    </Typography>
                  </CardContent>
                </Card>
              ))
            )}
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;
