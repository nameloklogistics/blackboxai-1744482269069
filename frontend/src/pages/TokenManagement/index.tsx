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
  TextField,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  CircularProgress,
} from '@mui/material';
import { AccountBalanceWallet, Send, Lock } from '@mui/icons-material';
import { useFormik } from 'formik';
import * as yup from 'yup';
import { useAppDispatch, useAppSelector } from '../../store';
import {
  fetchTokenBalance,
  fetchTransactionHistory,
  transferTokens,
  createEscrow,
} from '../../store/slices/tokenSlice';

const transferValidationSchema = yup.object({
  to: yup.string().required('Recipient address is required'),
  amount: yup
    .number()
    .required('Amount is required')
    .positive('Amount must be positive'),
  memo: yup.string(),
});

const TokenManagement: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const { balance, transactions, loading } = useAppSelector((state) => state.token);
  const [openTransferDialog, setOpenTransferDialog] = useState(false);

  useEffect(() => {
    if (user?.walletAddress) {
      dispatch(fetchTokenBalance(user.walletAddress));
      dispatch(fetchTransactionHistory(user.walletAddress));
    }
  }, [dispatch, user]);

  const transferFormik = useFormik({
    initialValues: {
      to: '',
      amount: '',
      memo: '',
    },
    validationSchema: transferValidationSchema,
    onSubmit: async (values) => {
      await dispatch(transferTokens(values));
      setOpenTransferDialog(false);
      transferFormik.resetForm();
    },
  });

  const formatDate = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  return (
    <Box sx={{ p: 3 }}>
      {/* Token Balance Card */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box display="flex" alignItems="center" mb={2}>
            <AccountBalanceWallet sx={{ fontSize: 40, mr: 2 }} color="primary" />
            <Box>
              <Typography variant="h6">LMT Balance</Typography>
              <Typography variant="h4" color="primary">
                {loading ? <CircularProgress size={24} /> : balance}
              </Typography>
            </Box>
          </Box>
          <Box display="flex" gap={2}>
            <Button
              variant="contained"
              startIcon={<Send />}
              onClick={() => setOpenTransferDialog(true)}
            >
              Transfer Tokens
            </Button>
            <Button variant="outlined" startIcon={<Lock />}>
              Create Escrow
            </Button>
          </Box>
        </CardContent>
      </Card>

      {/* Transaction History */}
      <Typography variant="h6" gutterBottom>
        Transaction History
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Type</TableCell>
              <TableCell>From</TableCell>
              <TableCell>To</TableCell>
              <TableCell>Amount</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Date</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <CircularProgress />
                </TableCell>
              </TableRow>
            ) : (
              transactions.map((tx) => (
                <TableRow key={tx.id}>
                  <TableCell>{tx.type}</TableCell>
                  <TableCell>{formatAddress(tx.from)}</TableCell>
                  <TableCell>{formatAddress(tx.to)}</TableCell>
                  <TableCell>{tx.amount} LMT</TableCell>
                  <TableCell>{tx.status}</TableCell>
                  <TableCell>{formatDate(tx.timestamp)}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Transfer Dialog */}
      <Dialog
        open={openTransferDialog}
        onClose={() => setOpenTransferDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <form onSubmit={transferFormik.handleSubmit}>
          <DialogTitle>Transfer Tokens</DialogTitle>
          <DialogContent>
            <Box sx={{ mt: 2 }}>
              <TextField
                fullWidth
                name="to"
                label="Recipient Address"
                value={transferFormik.values.to}
                onChange={transferFormik.handleChange}
                error={transferFormik.touched.to && Boolean(transferFormik.errors.to)}
                helperText={transferFormik.touched.to && transferFormik.errors.to}
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                name="amount"
                label="Amount (LMT)"
                type="number"
                value={transferFormik.values.amount}
                onChange={transferFormik.handleChange}
                error={transferFormik.touched.amount && Boolean(transferFormik.errors.amount)}
                helperText={transferFormik.touched.amount && transferFormik.errors.amount}
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                name="memo"
                label="Memo (Optional)"
                value={transferFormik.values.memo}
                onChange={transferFormik.handleChange}
                error={transferFormik.touched.memo && Boolean(transferFormik.errors.memo)}
                helperText={transferFormik.touched.memo && transferFormik.errors.memo}
              />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenTransferDialog(false)}>Cancel</Button>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={transferFormik.isSubmitting}
            >
              {transferFormik.isSubmitting ? <CircularProgress size={24} /> : 'Transfer'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default TokenManagement;
