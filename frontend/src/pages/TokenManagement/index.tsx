import React, { useEffect, useState } from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  useTheme,
} from '@mui/material';
import {
  AccountBalanceWallet,
  SwapHoriz,
  History,
  Lock,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../store';
import {
  fetchTokenBalance,
  transferTokens,
  fetchTokenTransactions,
} from '../../store/slices/tokenSlice';
import { useNotification } from '../../components/shared/NotificationCenter';
import PageHeader from '../../components/shared/PageHeader';
import DataTable from '../../components/shared/DataTable';
import { TokenTransaction } from '../../types';

const TokenManagement: React.FC = () => {
  const theme = useTheme();
  const dispatch = useAppDispatch();
  const { showSuccess, showError } = useNotification();
  const { user } = useAppSelector((state) => state.auth);
  const { balance, transactions, loading } = useAppSelector((state) => state.token);

  const [openTransferDialog, setOpenTransferDialog] = useState(false);
  const [transferAmount, setTransferAmount] = useState('');
  const [recipientAddress, setRecipientAddress] = useState('');
  const [transferLoading, setTransferLoading] = useState(false);

  useEffect(() => {
    if (user?.walletAddress) {
      dispatch(fetchTokenBalance(user.walletAddress));
      dispatch(fetchTokenTransactions(user.walletAddress));
    }
  }, [dispatch, user]);

  const handleTransferDialogOpen = () => {
    setOpenTransferDialog(true);
  };

  const handleTransferDialogClose = () => {
    setOpenTransferDialog(false);
    setTransferAmount('');
    setRecipientAddress('');
  };

  const handleTransfer = async () => {
    try {
      setTransferLoading(true);
      await dispatch(
        transferTokens({
          to: recipientAddress,
          amount: transferAmount,
          memo: 'Token transfer'
        })
      ).unwrap();
      showSuccess('Transfer completed successfully');
      handleTransferDialogClose();
      if (user?.walletAddress) {
        dispatch(fetchTokenBalance(user.walletAddress));
        dispatch(fetchTokenTransactions(user.walletAddress));
      }
    } catch (error) {
      showError('Transfer failed. Please try again.');
    } finally {
      setTransferLoading(false);
    }
  };

  const transactionColumns = [
    {
      id: 'timestamp',
      label: 'Date',
      format: (value: string) => new Date(value).toLocaleDateString(),
    },
    {
      id: 'type',
      label: 'Type',
      format: (value: string) => value.replace('_', ' '),
    },
    {
      id: 'amount',
      label: 'Amount',
      align: 'right' as const,
      format: (value: string) => `${value} LMT`,
    },
    {
      id: 'status',
      label: 'Status',
      align: 'center' as const,
    },
    {
      id: 'hash',
      label: 'Transaction Hash',
      format: (value: string) => `${value.slice(0, 8)}...${value.slice(-8)}`,
    },
  ];

  return (
    <Box>
      <PageHeader
        title="Token Management"
        subtitle="Manage your LMT tokens and view transaction history"
      />

      <Grid container spacing={3}>
        {/* Balance Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <AccountBalanceWallet
                  sx={{ fontSize: 40, color: theme.palette.primary.main, mr: 2 }}
                />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Available Balance
                  </Typography>
                  <Typography variant="h4">
                    {loading ? (
                      <CircularProgress size={24} />
                    ) : (
                      `${balance} LMT`
                    )}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Escrow Balance Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center">
                <Lock
                  sx={{ fontSize: 40, color: theme.palette.warning.main, mr: 2 }}
                />
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    In Escrow
                  </Typography>
                  <Typography variant="h4">0 LMT</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Actions */}
        <Grid item xs={12}>
          <Box display="flex" gap={2}>
            <Button
              variant="contained"
              startIcon={<SwapHoriz />}
              onClick={handleTransferDialogOpen}
            >
              Transfer Tokens
            </Button>
            <Button
              variant="outlined"
              startIcon={<History />}
            onClick={() => user?.walletAddress && dispatch(fetchTokenTransactions(user.walletAddress))}
            >
              Refresh History
            </Button>
          </Box>
        </Grid>

        {/* Transaction History */}
        <Grid item xs={12}>
          <DataTable<TokenTransaction>
            columns={transactionColumns}
            data={transactions}
            loading={loading}
            showSearch
            emptyStateMessage="No transactions found"
          />
        </Grid>
      </Grid>

      {/* Transfer Dialog */}
      <Dialog
        open={openTransferDialog}
        onClose={handleTransferDialogClose}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Transfer Tokens</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <TextField
              fullWidth
              label="Recipient Address"
              value={recipientAddress}
              onChange={(e) => setRecipientAddress(e.target.value)}
              margin="normal"
            />
            <TextField
              fullWidth
              label="Amount (LMT)"
              type="number"
              value={transferAmount}
              onChange={(e) => setTransferAmount(e.target.value)}
              margin="normal"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleTransferDialogClose}>Cancel</Button>
          <Button
            onClick={handleTransfer}
            variant="contained"
            disabled={
              transferLoading ||
              !transferAmount ||
              !recipientAddress ||
              parseFloat(transferAmount) <= 0
            }
          >
            {transferLoading ? <CircularProgress size={24} /> : 'Transfer'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default TokenManagement;
