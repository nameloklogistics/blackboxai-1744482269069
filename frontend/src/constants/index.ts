// API Configuration
export const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
export const API_TIMEOUT = 30000;

// Blockchain Configuration
export const BLOCKCHAIN_NETWORK = process.env.REACT_APP_STELLAR_NETWORK || 'TESTNET';
export const BLOCKCHAIN_HORIZON_URL = BLOCKCHAIN_NETWORK === 'TESTNET'
  ? 'https://horizon-testnet.stellar.org'
  : 'https://horizon.stellar.org';
export const TOKEN_CODE = 'LGST';
export const TOKEN_ISSUER = process.env.REACT_APP_TOKEN_ISSUER || '';

// Route Definitions
export const ROUTES = {
  LOGIN: '/login',
  REGISTER: '/register',
  DASHBOARD: '/dashboard',
  MARKETPLACE: {
    ROOT: '/marketplace',
    CREATE: '/marketplace/create',
    EDIT: '/marketplace/edit',
    MY_LISTINGS: '/marketplace/my-listings',
  },
  BOOKINGS: {
    ROOT: '/bookings',
    ACTIVE: '/bookings/active',
    COMPLETED: '/bookings/completed',
    DETAILS: '/bookings/:id',
  },
  TOKENS: {
    ROOT: '/tokens',
    WALLET: '/tokens/wallet',
    TRANSACTIONS: '/tokens/transactions',
    ESCROW: '/tokens/escrow',
  },
  PROFILE: {
    ROOT: '/profile',
    SETTINGS: '/profile/settings',
    SECURITY: '/profile/security',
  },
  ANALYTICS: {
    ROOT: '/analytics',
    OVERVIEW: '/analytics/overview',
    REPORTS: '/analytics/reports',
  },
};

// User Roles
export const USER_ROLES = {
  FREIGHT_FORWARDER: 'FREIGHT_FORWARDER',
  CUSTOMS_BROKER: 'CUSTOMS_BROKER',
  SHIPPER: 'SHIPPER',
  ADMIN: 'ADMIN',
} as const;

// Service Types
export const SERVICE_TYPES = {
  FREIGHT_FORWARDING: 'FREIGHT_FORWARDING',
  CUSTOMS_BROKERAGE: 'CUSTOMS_BROKERAGE',
  SHIPPING: 'SHIPPING',
  TRANSSHIPMENT: 'TRANSSHIPMENT',
} as const;

// Shipment Modes
export const SHIPMENT_MODES = {
  AIR: 'AIR',
  SEA: 'SEA',
  ROAD: 'ROAD',
  RAIL: 'RAIL',
} as const;

// Booking Statuses
export const BOOKING_STATUS = {
  PENDING: 'PENDING',
  CONFIRMED: 'CONFIRMED',
  IN_PROGRESS: 'IN_PROGRESS',
  COMPLETED: 'COMPLETED',
  CANCELLED: 'CANCELLED',
} as const;

// Payment Statuses
export const PAYMENT_STATUS = {
  UNPAID: 'UNPAID',
  PROCESSING: 'PROCESSING',
  PAID: 'PAID',
  REFUNDED: 'REFUNDED',
} as const;

// Transaction Types
export const TRANSACTION_TYPES = {
  PAYMENT: 'PAYMENT',
  ESCROW_CREATE: 'ESCROW_CREATE',
  ESCROW_RELEASE: 'ESCROW_RELEASE',
  REFUND: 'REFUND',
} as const;

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your internet connection.',
  SERVER_ERROR: 'Server error. Please try again later.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  INVALID_CREDENTIALS: 'Invalid email or password.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  TOKEN_ERROR: 'Error processing token transaction.',
  BOOKING_ERROR: 'Error processing booking.',
  LISTING_ERROR: 'Error processing service listing.',
} as const;

// Validation Rules
export const VALIDATION = {
  PASSWORD_MIN_LENGTH: 8,
  PASSWORD_REGEX: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  EMAIL_REGEX: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
  PHONE_REGEX: /^\+?[\d\s-]{10,}$/,
} as const;

// Date Formats
export const DATE_FORMATS = {
  DISPLAY: 'MMM DD, YYYY',
  DISPLAY_WITH_TIME: 'MMM DD, YYYY HH:mm',
  ISO: 'YYYY-MM-DD',
  ISO_WITH_TIME: 'YYYY-MM-DDTHH:mm:ss.SSSZ',
} as const;

// Pagination
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [10, 25, 50, 100],
} as const;

// File Upload
export const FILE_UPLOAD = {
  MAX_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_TYPES: ['image/jpeg', 'image/png', 'application/pdf'],
} as const;

// Local Storage Keys
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USER: 'user',
  THEME: 'theme',
  LANGUAGE: 'language',
} as const;

// Theme Settings
export const THEME = {
  DARK: 'dark',
  LIGHT: 'light',
} as const;

// Languages
export const LANGUAGES = {
  EN: 'en',
  ES: 'es',
  FR: 'fr',
} as const;

// Export all constants
export default {
  API_BASE_URL,
  API_TIMEOUT,
  BLOCKCHAIN_NETWORK,
  BLOCKCHAIN_HORIZON_URL,
  TOKEN_CODE,
  TOKEN_ISSUER,
  ROUTES,
  USER_ROLES,
  SERVICE_TYPES,
  SHIPMENT_MODES,
  BOOKING_STATUS,
  PAYMENT_STATUS,
  TRANSACTION_TYPES,
  ERROR_MESSAGES,
  VALIDATION,
  DATE_FORMATS,
  PAGINATION,
  FILE_UPLOAD,
  STORAGE_KEYS,
  THEME,
  LANGUAGES,
};
