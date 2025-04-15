import React from 'react';
import {
  Button,
  ButtonProps,
  CircularProgress,
  styled,
} from '@mui/material';

export interface LoadingButtonProps extends ButtonProps {
  loading?: boolean;
  loadingText?: string;
}

const StyledButton = styled(Button)(({ theme }) => ({
  position: 'relative',
  '& .MuiCircularProgress-root': {
    position: 'absolute',
    left: '50%',
    marginLeft: -12,
    marginTop: -12,
  },
  '&.Mui-disabled': {
    backgroundColor: theme.palette.action.disabledBackground,
  },
}));

const LoadingButton: React.FC<LoadingButtonProps> = ({
  children,
  loading = false,
  loadingText,
  disabled,
  startIcon,
  endIcon,
  ...props
}) => {
  return (
    <StyledButton
      disabled={loading || disabled}
      startIcon={loading ? undefined : startIcon}
      endIcon={loading ? undefined : endIcon}
      {...props}
    >
      {loading && <CircularProgress size={24} color="inherit" />}
      <span style={{ visibility: loading ? 'hidden' : 'visible' }}>
        {loading ? loadingText || children : children}
      </span>
    </StyledButton>
  );
};

export default LoadingButton;

// Example usage:
/*
import LoadingButton from '@components/shared/LoadingButton';

const MyComponent = () => {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await submitData();
    } finally {
      setLoading(false);
    }
  };

  return (
    <LoadingButton
      loading={loading}
      loadingText="Submitting..."
      variant="contained"
      color="primary"
      onClick={handleSubmit}
    >
      Submit
    </LoadingButton>
  );
};
*/
