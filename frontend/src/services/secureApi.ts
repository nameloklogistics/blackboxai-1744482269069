import { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { api } from './api';
import { setupSecurityMiddleware } from '../middleware/security';
import { 
  validation,
  loginValidationSchema,
  registerValidationSchema,
  tokenTransferValidationSchema,
  serviceListingValidationSchema,
  bookingValidationSchema
} from '../middleware/validation';
import crypto from '../utils/crypto';
import { 
  ApiResponse,
  User,
  TokenTransaction,
  ServiceListing,
  Booking,
  LoginFormData,
  RegisterFormData,
  ServiceListingFormData,
  BookingFormData
} from '../types';

class SecureApiService {
  private apiInstance: AxiosInstance;

  constructor() {
    this.apiInstance = api.client;
    this.setupSecurityMiddleware();
    this.setupRequestInterceptors();
    this.setupResponseInterceptors();
  }

  private setupSecurityMiddleware(): void {
    setupSecurityMiddleware(this.apiInstance);
  }

  private setupRequestInterceptors(): void {
    this.apiInstance.interceptors.request.use(
      async (config: InternalAxiosRequestConfig) => {
        try {
          // Validate and sanitize request data
          if (config.data) {
            // Sanitize input data
            config.data = validation.sanitizeObject(config.data);

            // Encrypt sensitive data
            config.data = crypto.encryptObject(config.data);
          }

          // Validate URL parameters
          if (config.params) {
            config.params = validation.sanitizeObject(config.params);
          }

          return config;
        } catch (error) {
          return Promise.reject(error);
        }
      },
      (error) => {
        return Promise.reject(error);
      }
    );
  }

  private setupResponseInterceptors(): void {
    this.apiInstance.interceptors.response.use(
      (response: AxiosResponse) => {
        try {
          // Validate response data structure
          if (!validation.validators.isValidJSON(JSON.stringify(response.data))) {
            throw new Error('Invalid response data structure');
          }

          // Decrypt sensitive data in response
          if (response.data) {
            response.data = crypto.decryptObject(response.data);
          }

          return response;
        } catch (error) {
          return Promise.reject(error);
        }
      },
      (error) => {
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints with enhanced security
  async login(formData: LoginFormData): Promise<ApiResponse<{ token: string; user: User }>> {
    // Validate input
    const validationResult = await validation.validateForm(
      loginValidationSchema,
      formData
    );

    if (!validationResult.success) {
      throw new Error(validationResult.errors?.[0] || 'Invalid input');
    }

    // Hash password before sending
    const hashedPassword = crypto.hash(formData.password);

    const response = await this.apiInstance.post('/auth/login', {
      email: formData.email,
      password: hashedPassword,
    });

    return response.data;
  }

  async register(formData: RegisterFormData): Promise<ApiResponse<{ token: string; user: User }>> {
    // Validate input
    const validationResult = await validation.validateForm(
      registerValidationSchema,
      formData
    );

    if (!validationResult.success) {
      throw new Error(validationResult.errors?.[0] || 'Invalid input');
    }

    // Hash password
    const hashedPassword = crypto.hash(formData.password);

    const response = await this.apiInstance.post('/auth/register', {
      ...formData,
      password: hashedPassword,
    });
    return response.data;
  }

  // Token endpoints with enhanced security
  async transferTokens(
    to: string,
    amount: string,
    memo?: string
  ): Promise<ApiResponse<TokenTransaction>> {
    // Validate input
    const validationResult = await validation.validateForm(
      tokenTransferValidationSchema,
      { to, amount, memo }
    );

    if (!validationResult.success) {
      throw new Error(validationResult.errors?.[0] || 'Invalid input');
    }

    const response = await this.apiInstance.post('/tokens/transfer', {
      to,
      amount,
      memo,
    });

    return response.data;
  }

  // Service listing endpoints with enhanced security
  async createServiceListing(
    formData: ServiceListingFormData
  ): Promise<ApiResponse<ServiceListing>> {
    // Validate input
    const validationResult = await validation.validateForm(
      serviceListingValidationSchema,
      formData
    );

    if (!validationResult.success) {
      throw new Error(validationResult.errors?.[0] || 'Invalid input');
    }

    const response = await this.apiInstance.post('/marketplace/listings', formData);
    return response.data;
  }

  // Booking endpoints with enhanced security
  async createBooking(formData: BookingFormData): Promise<ApiResponse<Booking>> {
    // Validate input
    const validationResult = await validation.validateForm(
      bookingValidationSchema,
      formData
    );

    if (!validationResult.success) {
      throw new Error(validationResult.errors?.[0] || 'Invalid input');
    }

    const response = await this.apiInstance.post('/bookings', formData);
    return response.data;
  }

  // Generic secure request method
  async request<T>(config: {
    method: string;
    url: string;
    data?: any;
    params?: any;
    schema?: any;
  }): Promise<ApiResponse<T>> {
    try {
      // Validate input if schema provided
      if (config.schema && config.data) {
        const validationResult = await validation.validateForm(
          config.schema,
          config.data
        );

        if (!validationResult.success) {
          throw new Error(validationResult.errors?.[0] || 'Invalid input');
        }
      }

      const response = await this.apiInstance.request({
        method: config.method,
        url: config.url,
        data: config.data,
        params: config.params,
      });

      return response.data;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Request failed: ${error.message}`);
      }
      throw new Error('Request failed');
    }
  }
}

export const secureApi = new SecureApiService();
export default secureApi;
