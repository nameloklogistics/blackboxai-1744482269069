import { format, parseISO } from 'date-fns';
import { TokenTransaction, BookingStatus, PaymentStatus, TransactionStatus } from '../types';
import { STORAGE_KEYS } from '../constants';

// Storage utilities
export const storage = {
  get: (key: keyof typeof STORAGE_KEYS): string | null => {
    try {
      return localStorage.getItem(STORAGE_KEYS[key]);
    } catch (error) {
      console.error(`Error getting item from storage: ${error}`);
      return null;
    }
  },

  set: (key: keyof typeof STORAGE_KEYS, value: string): void => {
    try {
      localStorage.setItem(STORAGE_KEYS[key], value);
    } catch (error) {
      console.error(`Error setting item in storage: ${error}`);
    }
  },

  remove: (key: keyof typeof STORAGE_KEYS): void => {
    try {
      localStorage.removeItem(STORAGE_KEYS[key]);
    } catch (error) {
      console.error(`Error removing item from storage: ${error}`);
    }
  },

  clear: (): void => {
    try {
      localStorage.clear();
    } catch (error) {
      console.error(`Error clearing storage: ${error}`);
    }
  },

  getJSON: <T>(key: keyof typeof STORAGE_KEYS): T | null => {
    try {
      const item = localStorage.getItem(STORAGE_KEYS[key]);
      return item ? JSON.parse(item) : null;
    } catch (error) {
      console.error(`Error getting JSON item from storage: ${error}`);
      return null;
    }
  },

  setJSON: <T>(key: keyof typeof STORAGE_KEYS, value: T): void => {
    try {
      localStorage.setItem(STORAGE_KEYS[key], JSON.stringify(value));
    } catch (error) {
      console.error(`Error setting JSON item in storage: ${error}`);
    }
  },

  exists: (key: keyof typeof STORAGE_KEYS): boolean => {
    try {
      return localStorage.getItem(STORAGE_KEYS[key]) !== null;
    } catch (error) {
      console.error(`Error checking item existence in storage: ${error}`);
      return false;
    }
  },

  isAvailable: (): boolean => {
    try {
      const test = '__storage_test__';
      localStorage.setItem(test, test);
      localStorage.removeItem(test);
      return true;
    } catch (error) {
      return false;
    }
  },
};


/**
 * Format a date string to a human-readable format
 */
export const formatDate = (dateString: string, formatString = 'MMM dd, yyyy'): string => {
  try {
    return format(parseISO(dateString), formatString);
  } catch (error) {
    console.error('Error formatting date:', error);
    return dateString;
  }
};

/**
 * Format a number as currency
 */
export const formatCurrency = (amount: number, currency = 'USD'): string => {
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(amount);
  } catch (error) {
    console.error('Error formatting currency:', error);
    return `${currency} ${amount}`;
  }
};

/**
 * Format a large number with abbreviations (K, M, B)
 */
export const formatNumber = (num: number): string => {
  const lookup = [
    { value: 1e9, symbol: 'B' },
    { value: 1e6, symbol: 'M' },
    { value: 1e3, symbol: 'K' },
  ];

  const rx = /\.0+$|(\.[0-9]*[1-9])0+$/;
  const item = lookup.find(item => num >= item.value);

  return item
    ? (num / item.value).toFixed(1).replace(rx, '$1') + item.symbol
    : num.toString();
};

/**
 * Truncate text with ellipsis
 */
export const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  return `${text.slice(0, maxLength)}...`;
};

/**
 * Format blockchain address for display
 */
export const formatAddress = (address: string, start = 6, end = 4): string => {
  if (!address) return '';
  if (address.length <= start + end) return address;
  return `${address.slice(0, start)}...${address.slice(-end)}`;
};

/**
 * Get status color for different status types
 */
export const getStatusColor = (
  status: BookingStatus | PaymentStatus | TransactionStatus
): string => {
  const statusColors = {
    // Booking statuses
    PENDING: '#ed6c02', // warning
    CONFIRMED: '#0288d1', // info
    IN_PROGRESS: '#1976d2', // primary
    COMPLETED: '#2e7d32', // success
    CANCELLED: '#d32f2f', // error
    // Payment statuses
    UNPAID: '#ed6c02', // warning
    PROCESSING: '#0288d1', // info
    PAID: '#2e7d32', // success
    REFUNDED: '#d32f2f', // error
    // Transaction statuses
    FAILED: '#d32f2f', // error
  };

  return statusColors[status] || '#64748b'; // neutral
};

/**
 * Format transaction type for display
 */
export const formatTransactionType = (type: TokenTransaction['type']): string => {
  const types = {
    PAYMENT: 'Payment',
    ESCROW_CREATE: 'Escrow Created',
    ESCROW_RELEASE: 'Escrow Released',
    REFUND: 'Refund',
  };

  return types[type] || type;
};

/**
 * Generate a random ID
 */
export const generateId = (): string => {
  return Math.random().toString(36).substring(2) + Date.now().toString(36);
};

/**
 * Deep clone an object
 */
export const deepClone = <T>(obj: T): T => {
  return JSON.parse(JSON.stringify(obj));
};

/**
 * Check if an object is empty
 */
export const isEmpty = (obj: object): boolean => {
  return Object.keys(obj).length === 0;
};

/**
 * Debounce a function
 */
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;

  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};

/**
 * Throttle a function
 */
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean;
  
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};

/**
 * Parse query string parameters
 */
export const parseQueryString = (queryString: string): Record<string, string> => {
  const params = new URLSearchParams(queryString);
  const result: Record<string, string> = {};
  
  params.forEach((value, key) => {
    result[key] = value;
  });
  
  return result;
};

/**
 * Convert object to query string
 */
export const objectToQueryString = (obj: Record<string, any>): string => {
  const params = new URLSearchParams();
  
  Object.entries(obj).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      params.append(key, String(value));
    }
  });
  
  return params.toString();
};

/**
 * Validate email address
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i;
  return emailRegex.test(email);
};

/**
 * Format file size
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};

/**
 * Get file extension
 */
export const getFileExtension = (filename: string): string => {
  return filename.slice((filename.lastIndexOf('.') - 1 >>> 0) + 2);
};

/**
 * Convert hex color to RGBA
 */
export const hexToRgba = (hex: string, alpha = 1): string => {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
};
