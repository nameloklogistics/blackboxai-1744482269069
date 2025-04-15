import React, { useState, useEffect } from 'react';
import {
  Snackbar,
  Alert,
  AlertTitle,
  Box,
  IconButton,
  Typography,
  Slide,
  SlideProps,
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../../store';

type NotificationType = 'success' | 'error' | 'warning' | 'info';

interface Notification {
  id: string;
  type: NotificationType;
  title?: string;
  message: string;
  autoHideDuration?: number;
}

const SlideTransition = (props: SlideProps) => {
  return <Slide {...props} direction="left" />;
};

const NotificationCenter: React.FC = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [openNotifications, setOpenNotifications] = useState<{ [key: string]: boolean }>({});

  // Add a new notification
  const addNotification = (notification: Omit<Notification, 'id'>) => {
    const id = Math.random().toString(36).substr(2, 9);
    const newNotification = {
      ...notification,
      id,
      autoHideDuration: notification.autoHideDuration || 6000,
    };

    setNotifications((prev) => [...prev, newNotification]);
    setOpenNotifications((prev) => ({ ...prev, [id]: true }));
  };

  // Remove a notification
  const removeNotification = (id: string) => {
    setNotifications((prev) => prev.filter((notification) => notification.id !== id));
    setOpenNotifications((prev) => {
      const newState = { ...prev };
      delete newState[id];
      return newState;
    });
  };

  // Handle notification close
  const handleClose = (id: string) => {
    setOpenNotifications((prev) => ({ ...prev, [id]: false }));
    setTimeout(() => removeNotification(id), 300); // Remove after exit animation
  };

  // Example usage of notification types
  const showSuccess = (message: string, title?: string) => {
    addNotification({
      type: 'success',
      title,
      message,
    });
  };

  const showError = (message: string, title?: string) => {
    addNotification({
      type: 'error',
      title,
      message,
      autoHideDuration: 10000, // Longer duration for errors
    });
  };

  const showWarning = (message: string, title?: string) => {
    addNotification({
      type: 'warning',
      title,
      message,
    });
  };

  const showInfo = (message: string, title?: string) => {
    addNotification({
      type: 'info',
      title,
      message,
    });
  };

  return (
    <Box
      sx={{
        position: 'fixed',
        top: (theme) => theme.spacing(2),
        right: (theme) => theme.spacing(2),
        zIndex: (theme) => theme.zIndex.snackbar,
      }}
    >
      {notifications.map((notification) => (
        <Snackbar
          key={notification.id}
          open={openNotifications[notification.id]}
          autoHideDuration={notification.autoHideDuration}
          onClose={() => handleClose(notification.id)}
          TransitionComponent={SlideTransition}
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
          sx={{ mb: 1 }}
        >
          <Alert
            severity={notification.type}
            variant="filled"
            sx={{ width: '100%', minWidth: '300px' }}
            action={
              <IconButton
                size="small"
                aria-label="close"
                color="inherit"
                onClick={() => handleClose(notification.id)}
              >
                <CloseIcon fontSize="small" />
              </IconButton>
            }
          >
            {notification.title && (
              <AlertTitle>{notification.title}</AlertTitle>
            )}
            <Typography variant="body2">{notification.message}</Typography>
          </Alert>
        </Snackbar>
      ))}
    </Box>
  );
};

// Create a notification context and hook for easy access
import { createContext, useContext } from 'react';

interface NotificationContextType {
  showSuccess: (message: string, title?: string) => void;
  showError: (message: string, title?: string) => void;
  showWarning: (message: string, title?: string) => void;
  showInfo: (message: string, title?: string) => void;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const notificationCenter = NotificationCenter();
  const value = {
    showSuccess: notificationCenter.showSuccess,
    showError: notificationCenter.showError,
    showWarning: notificationCenter.showWarning,
    showInfo: notificationCenter.showInfo,
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
      {notificationCenter}
    </NotificationContext.Provider>
  );
};

export const useNotification = () => {
  const context = useContext(NotificationContext);
  if (context === undefined) {
    throw new Error('useNotification must be used within a NotificationProvider');
  }
  return context;
};

// Example usage:
/*
import { useNotification } from './NotificationCenter';

const MyComponent = () => {
  const { showSuccess, showError } = useNotification();

  const handleSubmit = async () => {
    try {
      await submitData();
      showSuccess('Data submitted successfully!');
    } catch (error) {
      showError('Failed to submit data. Please try again.');
    }
  };

  return <button onClick={handleSubmit}>Submit</button>;
};
*/

export default NotificationCenter;
