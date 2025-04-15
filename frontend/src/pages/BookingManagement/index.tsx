import React, { useEffect, useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Grid,
  Typography,
  Stepper,
  Step,
  StepLabel,
  Chip,
  CircularProgress,
  Divider,
} from '@mui/material';
import {
  LocalShipping,
  Payment,
  Timeline,
  CheckCircle,
  Warning,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../store';
import {
  fetchBookings,
  updateBookingStatus,
  Booking,
} from '../../store/slices/marketplaceSlice';

const getStatusColor = (status: string) => {
  switch (status) {
    case 'CONFIRMED':
      return 'success';
    case 'PENDING':
      return 'warning';
    case 'IN_PROGRESS':
      return 'info';
    case 'COMPLETED':
      return 'success';
    case 'CANCELLED':
      return 'error';
    default:
      return 'default';
  }
};

const getPaymentStatusColor = (status: string) => {
  switch (status) {
    case 'PAID':
      return 'success';
    case 'UNPAID':
      return 'error';
    case 'PROCESSING':
      return 'warning';
    case 'REFUNDED':
      return 'info';
    default:
      return 'default';
  }
};

const BookingManagement: React.FC = () => {
  const dispatch = useAppDispatch();
  const { bookings, loading } = useAppSelector((state) => state.marketplace);
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);
  const [openTrackingDialog, setOpenTrackingDialog] = useState(false);

  useEffect(() => {
    dispatch(fetchBookings());
  }, [dispatch]);

  const handleViewTracking = (booking: Booking) => {
    setSelectedBooking(booking);
    setOpenTrackingDialog(true);
  };

  const getTrackingSteps = (trackingInfo: any[]) => {
    return trackingInfo.map((info) => ({
      label: info.status,
      description: info.description,
      location: info.location,
      timestamp: new Date(info.timestamp).toLocaleString(),
    }));
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Bookings
      </Typography>

      {loading ? (
        <Box display="flex" justifyContent="center" p={3}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={3}>
          {bookings.map((booking: Booking) => (
            <Grid item xs={12} md={6} key={booking.id}>
              <Card>
                <CardContent>
                  <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                    <Typography variant="h6">Booking #{booking.id}</Typography>
                    <Chip
                      label={booking.status}
                      color={getStatusColor(booking.status) as any}
                      size="small"
                    />
                  </Box>

                  <Typography color="textSecondary" gutterBottom>
                    Service ID: {booking.serviceId}
                  </Typography>

                  <Box mb={2}>
                    <Typography variant="subtitle2">Cargo Details:</Typography>
                    <Typography>Weight: {booking.cargoDetails.weight} kg</Typography>
                    <Typography>Volume: {booking.cargoDetails.volume} m³</Typography>
                    <Typography>Type: {booking.cargoDetails.type}</Typography>
                  </Box>

                  <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                    <Box display="flex" alignItems="center">
                      <Payment sx={{ mr: 1 }} />
                      <Typography>
                        Payment: {booking.paymentAmount} LMT
                      </Typography>
                    </Box>
                    <Chip
                      label={booking.paymentStatus}
                      color={getPaymentStatusColor(booking.paymentStatus) as any}
                      size="small"
                    />
                  </Box>

                  <Button
                    variant="outlined"
                    startIcon={<Timeline />}
                    fullWidth
                    onClick={() => handleViewTracking(booking)}
                  >
                    View Tracking
                  </Button>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Tracking Dialog */}
      <Dialog
        open={openTrackingDialog}
        onClose={() => setOpenTrackingDialog(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          Shipment Tracking - Booking #{selectedBooking?.id}
        </DialogTitle>
        <DialogContent>
          {selectedBooking && (
            <Box sx={{ mt: 2 }}>
              <Stepper orientation="vertical">
                {getTrackingSteps(selectedBooking.trackingInfo).map((step, index) => (
                  <Step key={index} active={true}>
                    <StepLabel
                      StepIconComponent={() =>
                        index === 0 ? (
                          <CheckCircle color="primary" />
                        ) : (
                          <LocalShipping color="primary" />
                        )
                      }
                    >
                      <Typography variant="subtitle1">{step.label}</Typography>
                      <Typography variant="body2" color="textSecondary">
                        {step.location}
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        {step.timestamp}
                      </Typography>
                    </StepLabel>
                  </Step>
                ))}
              </Stepper>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenTrackingDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default BookingManagement;
