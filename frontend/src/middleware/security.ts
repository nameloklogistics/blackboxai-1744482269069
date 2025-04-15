import { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosHeaders } from 'axios';
import { storage } from '../utils';
import { STORAGE_KEYS } from '../constants';

interface DevToolsDetector {
  isOpen: boolean;
  orientation: 'vertical' | 'horizontal' | null;
}

// Rate limiting configuration
const RATE_LIMIT = {
  maxRequests: 100,
  windowMs: 60000, // 1 minute
  requests: new Map<string, number[]>(),
};

// CSRF token management
let csrfToken: string | null = null;

/**
 * Generate a random CSRF token
 */
const generateCsrfToken = (): string => {
  return Math.random().toString(36).substr(2) + Date.now().toString(36);
};

/**
 * Check if the current time window has exceeded rate limit
 */
const isRateLimited = (endpoint: string): boolean => {
  const now = Date.now();
  const requests = RATE_LIMIT.requests.get(endpoint) || [];
  
  // Remove old requests outside the time window
  const validRequests = requests.filter(time => now - time < RATE_LIMIT.windowMs);
  RATE_LIMIT.requests.set(endpoint, validRequests);
  
  return validRequests.length >= RATE_LIMIT.maxRequests;
};

/**
 * Record a new request for rate limiting
 */
const recordRequest = (endpoint: string): void => {
  const requests = RATE_LIMIT.requests.get(endpoint) || [];
  requests.push(Date.now());
  RATE_LIMIT.requests.set(endpoint, requests);
};

/**
 * Setup security middleware for axios instance
 */
export const setupSecurityMiddleware = (axiosInstance: AxiosInstance): void => {
  // Request interceptor
  axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      if (!config.url) {
        return config;
      }

      // Rate limiting
      if (isRateLimited(config.url)) {
        throw new Error('Rate limit exceeded. Please try again later.');
      }
      recordRequest(config.url);

      // CSRF protection
      if (!csrfToken) {
        csrfToken = generateCsrfToken();
        if (csrfToken) {
          storage.set('TOKEN', csrfToken);
        }
      }

      // Add security headers
      const headers = new AxiosHeaders({
        'X-CSRF-Token': csrfToken || '',
        'X-Content-Type-Options': 'nosniff',
        'X-Frame-Options': 'DENY',
        'X-XSS-Protection': '1; mode=block',
        'Strict-Transport-Security': 'max-age=31536000; includeSubDomains',
      });

      config.headers = headers;

      // Encrypt sensitive data
      if (config.data && typeof config.data === 'object') {
        const sensitiveFields = ['password', 'secretKey', 'walletKey'];
        sensitiveFields.forEach(field => {
          if (config.data[field]) {
            config.data[field] = encryptSensitiveData(config.data[field]);
          }
        });
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor
  axiosInstance.interceptors.response.use(
    (response: AxiosResponse) => {
      // Validate response data
      if (!validateResponseData(response.data)) {
        throw new Error('Invalid response data received');
      }

      // Update CSRF token if provided
      const newCsrfToken = response.headers['x-csrf-token'];
      if (newCsrfToken && typeof newCsrfToken === 'string') {
        csrfToken = newCsrfToken;
        storage.set('TOKEN', newCsrfToken);
      }

      return response;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
};

/**
 * Encrypt sensitive data using AES encryption
 */
const encryptSensitiveData = (data: string): string => {
  // This is a placeholder for actual encryption implementation
  // In production, use a proper encryption library like crypto-js
  return btoa(data);
};

/**
 * Validate response data structure and types
 */
const validateResponseData = (data: any): boolean => {
  try {
    // Basic validation
    if (!data || typeof data !== 'object') {
      return false;
    }

    // Check for required fields based on response type
    if (data.hasOwnProperty('success')) {
      if (typeof data.success !== 'boolean') {
        return false;
      }
    }

    if (data.hasOwnProperty('data')) {
      // Additional validation based on expected data structure
      // This should be customized based on your API response format
    }

    return true;
  } catch (error) {
    console.error('Error validating response data:', error);
    return false;
  }
};

/**
 * Initialize security measures
 */
export const initializeSecurity = (): void => {
  // Prevent debugging
  preventDebugging();
  
  // Disable right-click
  disableRightClick();
  
  // Add visibility change detection
  handleVisibilityChange();
};

/**
 * Prevent debugging attempts
 */
const preventDebugging = (): void => {
  // Detect and prevent DevTools
  const devToolsDetector: DevToolsDetector = {
    isOpen: false,
    orientation: null
  };

  const emitDebuggerDetected = () => {
    window.dispatchEvent(new CustomEvent('debuggerDetected'));
  };

  // Detect DevTools by dimension changes
  setInterval(() => {
    const widthThreshold = window.outerWidth - window.innerWidth > 160;
    const heightThreshold = window.outerHeight - window.innerHeight > 160;
    
    if (widthThreshold || heightThreshold) {
      devToolsDetector.isOpen = true;
      devToolsDetector.orientation = widthThreshold ? 'vertical' : 'horizontal';
      emitDebuggerDetected();
    }
  }, 1000);

  // Prevent source map access
  if (process.env['NODE_ENV'] === 'production') {
    // @ts-ignore
    window.sourceMaps = undefined;
  }
};

/**
 * Disable right-click menu
 */
const disableRightClick = (): void => {
  document.addEventListener('contextmenu', (e) => {
    e.preventDefault();
    return false;
  });
};

/**
 * Handle visibility change events
 */
const handleVisibilityChange = (): void => {
  document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
      // Optionally clear sensitive data or take other security measures
      // when the application is not visible
    }
  });
};

// Export security utilities
export const security = {
  generateCsrfToken,
  validateResponseData,
  encryptSensitiveData,
};
