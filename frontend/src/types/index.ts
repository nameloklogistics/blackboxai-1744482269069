// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  errors?: string[];
}

export interface ApiError {
  status: number;
  message: string;
  errors?: string[];
}

// User Types
export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  walletAddress?: string;
  membershipStatus: MembershipStatus;
  createdAt: string;
  updatedAt: string;
}

export type UserRole = 
  | 'FREIGHT_FORWARDER'
  | 'CUSTOMS_BROKER'
  | 'SHIPPER'
  | 'ADMIN';

export type MembershipStatus = 
  | 'ACTIVE'
  | 'PENDING'
  | 'EXPIRED'
  | 'TRIAL';

// Service Types
export interface ServiceListing {
  id: string;
  providerId: string;
  providerName: string;
  serviceType: ServiceType;
  shipmentMode: ShipmentMode;
  origin: string;
  destination: string;
  rate: number;
  description: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export type ServiceType =
  | 'FREIGHT_FORWARDING'
  | 'CUSTOMS_BROKERAGE'
  | 'SHIPPING'
  | 'TRANSSHIPMENT';

export type ShipmentMode =
  | 'AIR'
  | 'SEA'
  | 'ROAD'
  | 'RAIL';

// Booking Types
export interface Booking {
  id: string;
  serviceId: string;
  customerId: string;
  providerId: string;
  status: BookingStatus;
  cargoDetails: CargoDetails;
  paymentStatus: PaymentStatus;
  paymentAmount: number;
  trackingInfo: TrackingInfo[];
  createdAt: string;
  updatedAt: string;
}

export type BookingStatus =
  | 'PENDING'
  | 'CONFIRMED'
  | 'IN_PROGRESS'
  | 'COMPLETED'
  | 'CANCELLED';

export type PaymentStatus =
  | 'UNPAID'
  | 'PROCESSING'
  | 'PAID'
  | 'REFUNDED';

export interface CargoDetails {
  weight: number;
  volume: number;
  type: string;
  description: string;
  hazardous: boolean;
  specialInstructions?: string;
}

export interface TrackingInfo {
  timestamp: string;
  status: string;
  location: string;
  description: string;
}

// Token Types
export interface TokenTransaction {
  id: string;
  type: TransactionType;
  from: string;
  to: string;
  amount: string;
  status: TransactionStatus;
  timestamp: string;
  memo?: string;
}

export type TransactionType =
  | 'PAYMENT'
  | 'ESCROW_CREATE'
  | 'ESCROW_RELEASE'
  | 'REFUND';

export type TransactionStatus =
  | 'PENDING'
  | 'COMPLETED'
  | 'FAILED';

export interface EscrowBalance {
  id: string;
  bookingId: string;
  amount: string;
  createdAt: string;
  expiresAt: string;
  status: EscrowStatus;
}

export type EscrowStatus =
  | 'ACTIVE'
  | 'RELEASED'
  | 'EXPIRED'
  | 'REFUNDED';

// Form Types
export interface LoginFormData {
  email: string;
  password: string;
}

export interface RegisterFormData {
  email: string;
  password: string;
  name: string;
  role: UserRole;
}

export interface ServiceListingFormData {
  serviceType: ServiceType;
  shipmentMode: ShipmentMode;
  origin: string;
  destination: string;
  rate: number;
  description: string;
}

export interface BookingFormData {
  serviceId: string;
  cargoDetails: CargoDetails;
}

// Component Props Types
export interface LoadingButtonProps {
  loading?: boolean;
  children: React.ReactNode;
  onClick?: () => void;
  [key: string]: any;
}

export interface StatusChipProps {
  status: BookingStatus | PaymentStatus | TransactionStatus | EscrowStatus;
  [key: string]: any;
}

export interface PageHeaderProps {
  title: string;
  subtitle?: string;
  actions?: React.ReactNode;
}

export interface DataTableProps<T> {
  data: T[];
  columns: TableColumn[];
  loading?: boolean;
  onRowClick?: (row: T) => void;
  pagination?: boolean;
  [key: string]: any;
}

export interface TableColumn {
  id: string;
  label: string;
  render?: (value: any, row: any) => React.ReactNode;
  sortable?: boolean;
  width?: string | number;
}

export interface ConfirmationDialogProps {
  open: boolean;
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  confirmText?: string;
  cancelText?: string;
}

// Store Types
export interface RootState {
  auth: AuthState;
  marketplace: MarketplaceState;
  token: TokenState;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

export interface MarketplaceState {
  listings: ServiceListing[];
  bookings: Booking[];
  selectedListing: ServiceListing | null;
  selectedBooking: Booking | null;
  loading: boolean;
  error: string | null;
  filters: MarketplaceFilters;
}

export interface MarketplaceFilters extends SearchFilter {
  minRate?: number;
  maxRate?: number;
}

export interface SearchFilter {
  page: number;
  pageSize: number;
}

export interface MarketplaceFilters extends SearchFilter {
  minRate?: number;
  maxRate?: number;
  query?: string;
  serviceType?: ServiceType;
  shipmentMode?: ShipmentMode;
  origin?: string;
  destination?: string;
  minPrice?: number;
  maxPrice?: number;
  routeType?: string;
  transitTime?: string;
  providerId?: string;
  isActive?: boolean;
  validFrom?: string;
  validUntil?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface SearchResult {
  services: ServiceListing[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  totalPages: number;
}

export interface SearchSuggestion {
  type: 'location' | 'service' | 'provider';
  value: string;
  label: string;
}

export interface TokenState {
  balance: string;
  transactions: TokenTransaction[];
  escrowBalances: EscrowBalance[];
  loading: boolean;
  error: string | null;
}
