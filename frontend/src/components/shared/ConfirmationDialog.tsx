import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  Box,
  Typography,
  IconButton,
} from '@mui/material';
import {
  Warning,
  Error,
  Info,
  CheckCircle,
  Close as CloseIcon,
} from '@mui/icons-material';

type DialogType = 'warning' | 'error' | 'info' | 'success';

interface ConfirmationDialogProps {
  open: boolean;
  title: string;
  message: string | React.ReactNode;
  type?: DialogType;
  confirmLabel?: string;
  cancelLabel?: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
  maxWidth?: 'xs' | 'sm' | 'md';
}

const getDialogIcon = (type: DialogType) => {
  switch (type) {
    case 'warning':
      return <Warning color="warning" sx={{ fontSize: 40 }} />;
    case 'error':
      return <Error color="error" sx={{ fontSize: 40 }} />;
    case 'success':
      return <CheckCircle color="success" sx={{ fontSize: 40 }} />;
    default:
      return <Info color="info" sx={{ fontSize: 40 }} />;
  }
};

const getDialogColor = (type: DialogType) => {
  switch (type) {
    case 'warning':
      return 'warning.main';
    case 'error':
      return 'error.main';
    case 'success':
      return 'success.main';
    default:
      return 'info.main';
  }
};

const ConfirmationDialog: React.FC<ConfirmationDialogProps> = ({
  open,
  title,
  message,
  type = 'warning',
  confirmLabel = 'Confirm',
  cancelLabel = 'Cancel',
  onConfirm,
  onCancel,
  loading = false,
  maxWidth = 'sm',
}) => {
  return (
    <Dialog
      open={open}
      onClose={onCancel}
      maxWidth={maxWidth}
      fullWidth
      aria-labelledby="confirmation-dialog-title"
    >
      <DialogTitle id="confirmation-dialog-title" sx={{ m: 0, p: 2 }}>
        <Box display="flex" alignItems="center">
          <Box sx={{ mr: 1 }}>{getDialogIcon(type)}</Box>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            {title}
          </Typography>
          <IconButton
            aria-label="close"
            onClick={onCancel}
            sx={{
              position: 'absolute',
              right: 8,
              top: 8,
              color: (theme) => theme.palette.grey[500],
            }}
          >
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {typeof message === 'string' ? (
          <DialogContentText>{message}</DialogContentText>
        ) : (
          message
        )}
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button
          onClick={onCancel}
          color="inherit"
          disabled={loading}
          variant="outlined"
        >
          {cancelLabel}
        </Button>
        <Button
          onClick={onConfirm}
          color={type}
          variant="contained"
          disabled={loading}
          autoFocus
        >
          {confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

// Predefined confirmation dialogs for common use cases
export const DeleteConfirmationDialog: React.FC<{
  open: boolean;
  itemName: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}> = ({ open, itemName, onConfirm, onCancel, loading }) => (
  <ConfirmationDialog
    open={open}
    title="Confirm Deletion"
    message={`Are you sure you want to delete ${itemName}? This action cannot be undone.`}
    type="error"
    confirmLabel="Delete"
    onConfirm={onConfirm}
    onCancel={onCancel}
    loading={loading}
  />
);

export const CancelBookingDialog: React.FC<{
  open: boolean;
  bookingId: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}> = ({ open, bookingId, onConfirm, onCancel, loading }) => (
  <ConfirmationDialog
    open={open}
    title="Cancel Booking"
    message={`Are you sure you want to cancel booking #${bookingId}? This action cannot be undone.`}
    type="warning"
    confirmLabel="Cancel Booking"
    onConfirm={onConfirm}
    onCancel={onCancel}
    loading={loading}
  />
);

export const ProcessPaymentDialog: React.FC<{
  open: boolean;
  amount: number;
  currency: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}> = ({ open, amount, currency, onConfirm, onCancel, loading }) => (
  <ConfirmationDialog
    open={open}
    title="Confirm Payment"
    message={`Are you sure you want to process payment of ${amount} ${currency}?`}
    type="info"
    confirmLabel="Process Payment"
    onConfirm={onConfirm}
    onCancel={onCancel}
    loading={loading}
  />
);

export default ConfirmationDialog;
