import React from 'react';
import { Chip, ChipProps } from '@mui/material';
import {
  CheckCircle,
  Warning,
  Error,
  Info,
  Pending,
  Block,
} from '@mui/icons-material';

type StatusType =
  | 'success'
  | 'warning'
  | 'error'
  | 'info'
  | 'pending'
  | 'inactive';

interface StatusChipProps extends Omit<ChipProps, 'color'> {
  status: StatusType;
  text?: string;
}

const getStatusConfig = (status: StatusType) => {
  switch (status) {
    case 'success':
      return {
        color: 'success' as const,
        icon: CheckCircle,
        defaultText: 'Success',
      };
    case 'warning':
      return {
        color: 'warning' as const,
        icon: Warning,
        defaultText: 'Warning',
      };
    case 'error':
      return {
        color: 'error' as const,
        icon: Error,
        defaultText: 'Error',
      };
    case 'info':
      return {
        color: 'info' as const,
        icon: Info,
        defaultText: 'Info',
      };
    case 'pending':
      return {
        color: 'warning' as const,
        icon: Pending,
        defaultText: 'Pending',
      };
    case 'inactive':
      return {
        color: 'default' as const,
        icon: Block,
        defaultText: 'Inactive',
      };
    default:
      return {
        color: 'default' as const,
        icon: Info,
        defaultText: 'Unknown',
      };
  }
};

const StatusChip: React.FC<StatusChipProps> = ({
  status,
  text,
  size = 'small',
  ...props
}) => {
  const { color, icon: Icon, defaultText } = getStatusConfig(status);

  return (
    <Chip
      icon={<Icon sx={{ fontSize: size === 'small' ? 16 : 20 }} />}
      label={text || defaultText}
      color={color}
      size={size}
      {...props}
    />
  );
};

export const BookingStatusChip: React.FC<{ status: string }> = ({ status }) => {
  const getBookingStatus = (): StatusType => {
    switch (status) {
      case 'CONFIRMED':
        return 'success';
      case 'PENDING':
        return 'pending';
      case 'IN_PROGRESS':
        return 'info';
      case 'COMPLETED':
        return 'success';
      case 'CANCELLED':
        return 'error';
      default:
        return 'info';
    }
  };

  return <StatusChip status={getBookingStatus()} text={status} />;
};

export const PaymentStatusChip: React.FC<{ status: string }> = ({ status }) => {
  const getPaymentStatus = (): StatusType => {
    switch (status) {
      case 'PAID':
        return 'success';
      case 'UNPAID':
        return 'error';
      case 'PROCESSING':
        return 'pending';
      case 'REFUNDED':
        return 'info';
      default:
        return 'info';
    }
  };

  return <StatusChip status={getPaymentStatus()} text={status} />;
};

export const MembershipStatusChip: React.FC<{ status: string }> = ({ status }) => {
  const getMembershipStatus = (): StatusType => {
    switch (status) {
      case 'ACTIVE':
        return 'success';
      case 'TRIAL':
        return 'info';
      case 'EXPIRED':
        return 'error';
      case 'PENDING':
        return 'pending';
      default:
        return 'inactive';
    }
  };

  return <StatusChip status={getMembershipStatus()} text={status} />;
};

export default StatusChip;
