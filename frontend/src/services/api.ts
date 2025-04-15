import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios';
import { ApiResponse, ApiError, User, ServiceListing, Booking, TokenTransaction } from '../types';
import { API_BASE_URL, ERROR_MESSAGES } from '../constants';

class ApiService {
  public client: AxiosInstance;

  public createClient(): AxiosInstance {
    return axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  constructor() {
    this.client = this.createClient();
    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(this.handleError(error));
      }
    );
  }

  private handleError(error: AxiosError): ApiError {
    if (error.response) {
      return {
        status: error.response.status,
        message: (error.response.data as any).message || ERROR_MESSAGES.SERVER_ERROR,
        errors: (error.response.data as any).errors,
      };
    }
    if (error.request) {
      return {
        status: 0,
        message: ERROR_MESSAGES.NETWORK_ERROR,
      };
    }
    return {
      status: 500,
      message: ERROR_MESSAGES.SERVER_ERROR,
    };
  }

  // Auth endpoints
  async login(email: string, password: string): Promise<ApiResponse<{ token: string; user: User }>> {
    const response = await this.client.post('/auth/login', { email, password });
    return response.data;
  }

  async register(userData: Partial<User>): Promise<ApiResponse<{ token: string; user: User }>> {
    const response = await this.client.post('/auth/register', userData);
    return response.data;
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    const response = await this.client.get('/auth/me');
    return response.data;
  }

  // Service Listings endpoints
  async getServiceListings(): Promise<ApiResponse<ServiceListing[]>> {
    const response = await this.client.get('/marketplace/listings');
    return response.data;
  }

  async createServiceListing(listing: Partial<ServiceListing>): Promise<ApiResponse<ServiceListing>> {
    const response = await this.client.post('/marketplace/listings', listing);
    return response.data;
  }

  async updateServiceListing(
    id: string,
    updates: Partial<ServiceListing>
  ): Promise<ApiResponse<ServiceListing>> {
    const response = await this.client.put(`/marketplace/listings/${id}`, updates);
    return response.data;
  }

  async deleteServiceListing(id: string): Promise<void> {
    await this.client.delete(`/marketplace/listings/${id}`);
  }

  // Booking endpoints
  async getBookings(): Promise<ApiResponse<Booking[]>> {
    const response = await this.client.get('/bookings');
    return response.data;
  }

  async createBooking(booking: Partial<Booking>): Promise<ApiResponse<Booking>> {
    const response = await this.client.post('/bookings', booking);
    return response.data;
  }

  async updateBookingStatus(id: string, status: string): Promise<ApiResponse<Booking>> {
    const response = await this.client.put(`/bookings/${id}/status`, { status });
    return response.data;
  }

  // Token endpoints
  async getTokenBalance(address: string): Promise<ApiResponse<string>> {
    const response = await this.client.get(`/tokens/balance/${address}`);
    return response.data;
  }

  async getTokenTransactions(address: string): Promise<ApiResponse<TokenTransaction[]>> {
    const response = await this.client.get(`/tokens/transactions/${address}`);
    return response.data;
  }

  async transferTokens(
    to: string,
    amount: string,
    memo?: string
  ): Promise<ApiResponse<TokenTransaction>> {
    const response = await this.client.post('/tokens/transfer', { to, amount, memo });
    return response.data;
  }

  async createEscrow(
    amount: string,
    bookingId: string,
    duration: number
  ): Promise<ApiResponse<TokenTransaction>> {
    const response = await this.client.post('/tokens/escrow', { amount, bookingId, duration });
    return response.data;
  }

  async releaseEscrow(escrowId: string): Promise<ApiResponse<TokenTransaction>> {
    const response = await this.client.post(`/tokens/escrow/${escrowId}/release`);
    return response.data;
  }

  // Wallet endpoints
  async connectWallet(secretKey: string, signature: string): Promise<ApiResponse<User>> {
    const response = await this.client.post('/wallet/connect', { secretKey, signature });
    return response.data;
  }

  async disconnectWallet(): Promise<ApiResponse<User>> {
    const response = await this.client.post('/wallet/disconnect');
    return response.data;
  }

  // Profile endpoints
  async updateProfile(updates: Partial<User>): Promise<ApiResponse<User>> {
    const response = await this.client.put('/profile', updates);
    return response.data;
  }
}

export const api = new ApiService();
export default api;
