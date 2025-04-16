import React, { useEffect, useState } from 'react';
import {
  Box,
  Tabs,
  Tab,
  Typography,
  Card,
  CardContent,
  Grid,
  Chip,
  Button,
  CircularProgress,
  useTheme,
} from '@mui/material';
import {
  LocalShipping,
  AccessTime,
  CheckCircle,
  Cancel,
  Payment,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../store';
import { fetchBookings } from '../../store/slices/marketplaceSlice';
import PageHeader from '../../components/shared/PageHeader';
import { BookingStatusChip } from '../../components/shared/StatusChip';
import DataTable from '../../components/shared/DataTable';
import { Booking, BookingStatus } from '../../types';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const TabPanel = (props: TabPanelProps) => {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`booking-tabpanel-${index}`}
      aria-labelledby={`booking-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
};

const BookingManagement: React.FC = () => {
  const theme = useTheme();
  const dispatch = useAppDispatch();
  const [tabValue, setTabValue] = useState(0);
  const { bookings, loading } = useAppSelector((state) => state.marketplace);

  useEffect(() => {
    dispatch(fetchBookings());
  }, [dispatch]);

  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const getStatusColor = (status: BookingStatus) => {
    switch (status) {
      case 'PENDING':
        return theme.palette.warning.main;
      case 'CONFIRMED':
        return theme.palette.info.main;
      case 'IN_PROGRESS':
        return theme.palette.primary.main;
      case 'COMPLETED':
        return theme.palette.success.main;
      case 'CANCELLED':
        return theme.palette.error.main;
      default:
        return theme.palette.grey[500];
    }
  };

  const columns = [
    {
      id: 'id',
      label: 'Booking ID',
      render: (value: string) => (
        <Typography variant="body2" color="primary">
          #{value}
        </Typography>
      ),
    },
    {
      id: 'status',
      label: 'Status',
      render: (value: BookingStatus) => (
        <BookingStatusChip status={value} />
      ),
    },
    {
      id: 'paymentStatus',
      label: 'Payment',
      render: (value: string) => (
        <Chip
          size="small"
          label={value}
          color={value === 'PAID' ? 'success' : 'warning'}
        />
      ),
    },
    {
      id: 'createdAt',
      label: 'Created',
      render: (value: string) => new Date(value).toLocaleDateString(),
    },
    {
      id: 'paymentAmount',
      label: 'Amount',
      render: (value: number) => `${value} LMT`,
    },
  ];

  const filterBookings = (status: BookingStatus[]) => {
    return bookings.filter((booking) => status.includes(booking.status));
  };

  return (
    <Box>
      <PageHeader
        title="Booking Management"
        subtitle="Track and manage your logistics bookings"
      />

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="All Bookings" />
          <Tab label="Active" />
          <Tab label="Completed" />
          <Tab label="Cancelled" />
        </Tabs>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" p={3}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <TabPanel value={tabValue} index={0}>
            <DataTable
              data={bookings}
              columns={columns}
              onRowClick={(booking: Booking) => console.log('Clicked booking:', booking)}
            />
          </TabPanel>

          <TabPanel value={tabValue} index={1}>
            <DataTable
              data={filterBookings(['CONFIRMED', 'IN_PROGRESS'])}
              columns={columns}
              onRowClick={(booking: Booking) => console.log('Clicked booking:', booking)}
            />
          </TabPanel>

          <TabPanel value={tabValue} index={2}>
            <DataTable
              data={filterBookings(['COMPLETED'])}
              columns={columns}
              onRowClick={(booking: Booking) => console.log('Clicked booking:', booking)}
            />
          </TabPanel>

          <TabPanel value={tabValue} index={3}>
            <DataTable
              data={filterBookings(['CANCELLED'])}
              columns={columns}
              onRowClick={(booking: Booking) => console.log('Clicked booking:', booking)}
            />
          </TabPanel>
        </>
      )}

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mt: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <LocalShipping color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Total Bookings
                  </Typography>
                  <Typography variant="h6">{bookings.length}</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <AccessTime color="warning" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Pending
                  </Typography>
                  <Typography variant="h6">
                    {filterBookings(['PENDING']).length}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <CheckCircle color="success" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Completed
                  </Typography>
                  <Typography variant="h6">
                    {filterBookings(['COMPLETED']).length}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <Payment color="info" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Total Value
                  </Typography>
                  <Typography variant="h6">
                    {bookings.reduce((sum, booking) => sum + booking.paymentAmount, 0)} LMT
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default BookingManagement;
