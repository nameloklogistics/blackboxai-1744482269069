import React, { useEffect, useState, useCallback } from 'react';
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
  TextField,
  Typography,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  CircularProgress,
  Pagination,
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { useFormik } from 'formik';
import * as yup from 'yup';
import { useAppDispatch, useAppSelector } from '../../store';
import {
  fetchServiceListings,
  createServiceListing,
  searchServices,
  setFilters,
  clearFilters,
  selectSearchResults,
} from '../../store/slices/marketplaceSlice';
import { SearchBar } from '../../components/shared/SearchBar';
import { ServiceListing, ServiceType, ShipmentMode } from '../../types';

interface ServiceListingFormData {
  serviceType: ServiceType;
  shipmentMode: ShipmentMode;
  origin: string;
  destination: string;
  rate: string;
  description: string;
}

const validationSchema = yup.object({
  serviceType: yup.string().required('Service type is required'),
  shipmentMode: yup.string().required('Shipment mode is required'),
  origin: yup.string().required('Origin is required'),
  destination: yup.string().required('Destination is required'),
  rate: yup.number().required('Rate is required').positive('Rate must be positive'),
  description: yup.string().required('Description is required'),
});

const ServiceListings: React.FC = () => {
  const dispatch = useAppDispatch();
  const { listings, loading, filters } = useAppSelector((state) => state.marketplace);
  const searchResults = useAppSelector(selectSearchResults);
  const [openDialog, setOpenDialog] = useState(false);

  const loadListings = useCallback(async () => {
    try {
      if (!filters.query) {
        const resultAction = await dispatch(fetchServiceListings());
        if (fetchServiceListings.fulfilled.match(resultAction)) {
          // Handle success if needed
        }
      } else {
        const resultAction = await dispatch(searchServices(filters));
        if (searchServices.fulfilled.match(resultAction)) {
          // Handle success if needed
        }
      }
    } catch (error) {
      console.error('Failed to fetch listings:', error);
    }
  }, [dispatch, filters]);

  useEffect(() => {
    loadListings();
  }, [filters]); // Re-run loadListings when filters change

  const handleSearch = (query: string) => {
    dispatch(setFilters({ ...filters, query, page: 1 }));
  };

  const handlePageChange = (_: React.ChangeEvent<unknown>, page: number) => {
    if (searchResults) {
      dispatch(setFilters({ ...filters, page }));
    }
  };

  const formik = useFormik<ServiceListingFormData>({
    initialValues: {
      serviceType: 'FREIGHT_FORWARDING',
      shipmentMode: 'AIR',
      origin: '',
      destination: '',
      rate: '',
      description: '',
    },
    validationSchema,
    onSubmit: async (values) => {
      try {
        const createAction = await dispatch(createServiceListing({
          ...values,
          rate: parseFloat(values.rate),
        }));
        
        if (createServiceListing.fulfilled.match(createAction)) {
          setOpenDialog(false);
          formik.resetForm();
          loadListings(); // Use loadListings instead of separate dispatch calls
        }
      } catch (error) {
        console.error('Failed to create listing:', error);
      }
    },
  });

  const handleOpenDialog = () => {
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    formik.resetForm();
  };

  return (
    <Box sx={{ p: 3 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Service Listings</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={handleOpenDialog}
        >
          Create New Listing
        </Button>
      </Box>

      <Box mb={3}>
        <SearchBar
          onSearch={handleSearch}
          placeholder="Search services by name, location, or provider..."
          fullWidth
        />
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" p={3}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Grid container spacing={3}>
            {listings.map((listing: ServiceListing) => (
              <Grid item xs={12} md={6} lg={4} key={listing.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {listing.serviceType}
                    </Typography>
                    <Typography color="textSecondary" gutterBottom>
                      {listing.shipmentMode}
                    </Typography>
                    <Typography variant="body1" gutterBottom>
                      From: {listing.origin}
                    </Typography>
                    <Typography variant="body1" gutterBottom>
                      To: {listing.destination}
                    </Typography>
                    <Typography variant="h5" color="primary" gutterBottom>
                      {listing.rate} LMT
                    </Typography>
                    <Typography variant="body2">{listing.description}</Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>

          {searchResults && (
            <Box display="flex" justifyContent="center" mt={3}>
              <Pagination
                count={searchResults.totalPages}
                page={searchResults.currentPage}
                onChange={handlePageChange}
                color="primary"
              />
            </Box>
          )}
        </>
      )}

      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <form onSubmit={formik.handleSubmit}>
          <DialogTitle>Create New Service Listing</DialogTitle>
          <DialogContent>
            <Box sx={{ mt: 2 }}>
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Service Type</InputLabel>
                <Select
                  name="serviceType"
                  value={formik.values.serviceType}
                  onChange={formik.handleChange}
                  error={Boolean(formik.touched.serviceType && formik.errors.serviceType)}
                  label="Service Type"
                >
                  <MenuItem value="FREIGHT_FORWARDING">Freight Forwarding</MenuItem>
                  <MenuItem value="CUSTOMS_BROKERAGE">Customs Brokerage</MenuItem>
                  <MenuItem value="SHIPPING">Shipping</MenuItem>
                  <MenuItem value="TRANSSHIPMENT">Transshipment</MenuItem>
                </Select>
              </FormControl>

              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Shipment Mode</InputLabel>
                <Select
                  name="shipmentMode"
                  value={formik.values.shipmentMode}
                  onChange={formik.handleChange}
                  error={Boolean(formik.touched.shipmentMode && formik.errors.shipmentMode)}
                  label="Shipment Mode"
                >
                  <MenuItem value="AIR">Air</MenuItem>
                  <MenuItem value="SEA">Sea</MenuItem>
                  <MenuItem value="ROAD">Road</MenuItem>
                  <MenuItem value="RAIL">Rail</MenuItem>
                </Select>
              </FormControl>

              <TextField
                fullWidth
                name="origin"
                label="Origin"
                value={formik.values.origin}
                onChange={formik.handleChange}
                error={Boolean(formik.touched.origin && formik.errors.origin)}
                helperText={formik.touched.origin && formik.errors.origin}
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                name="destination"
                label="Destination"
                value={formik.values.destination}
                onChange={formik.handleChange}
                error={Boolean(formik.touched.destination && formik.errors.destination)}
                helperText={formik.touched.destination && formik.errors.destination}
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                name="rate"
                label="Rate (LMT)"
                type="number"
                value={formik.values.rate}
                onChange={formik.handleChange}
                error={Boolean(formik.touched.rate && formik.errors.rate)}
                helperText={formik.touched.rate && formik.errors.rate}
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                name="description"
                label="Description"
                multiline
                rows={4}
                value={formik.values.description}
                onChange={formik.handleChange}
                error={Boolean(formik.touched.description && formik.errors.description)}
                helperText={formik.touched.description && formik.errors.description}
              />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseDialog}>Cancel</Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={formik.isSubmitting}
            >
              {formik.isSubmitting ? <CircularProgress size={24} /> : 'Create'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default ServiceListings;
