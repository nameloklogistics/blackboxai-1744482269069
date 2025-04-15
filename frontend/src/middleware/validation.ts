import { z } from 'zod';
import { VALIDATION } from '../constants';

// Base validation schemas
const emailSchema = z
  .string()
  .email('Invalid email format')
  .regex(VALIDATION.EMAIL_REGEX, 'Invalid email format');

const passwordSchema = z
  .string()
  .min(VALIDATION.PASSWORD_MIN_LENGTH, `Password must be at least ${VALIDATION.PASSWORD_MIN_LENGTH} characters`)
  .regex(
    VALIDATION.PASSWORD_REGEX,
    'Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character'
  );

const walletAddressSchema = z
  .string()
  .length(56, 'Invalid wallet address length')
  .regex(/^[A-Z0-9]+$/, 'Invalid wallet address format');

const amountSchema = z
  .string()
  .regex(/^\d+(\.\d{0,7})?$/, 'Invalid amount format')
  .refine((val: string) => parseFloat(val) > 0, 'Amount must be greater than 0');

// Form validation schemas
export const loginValidationSchema = z.object({
  email: emailSchema,
  password: passwordSchema,
});

export const registerValidationSchema = z.object({
  email: emailSchema.optional(),
  password: passwordSchema.optional(),
  name: z.string().min(2, 'Name is too short').optional(),
  role: z.enum(['FREIGHT_FORWARDER', 'CUSTOMS_BROKER', 'SHIPPER', 'ADMIN']).optional(),
}).refine((data) => {
  // At least one field must be present
  return Object.keys(data).length > 0;
}, 'At least one field must be provided');

export const tokenTransferValidationSchema = z.object({
  to: walletAddressSchema,
  amount: amountSchema,
  memo: z.string().optional(),
});

export const serviceListingValidationSchema = z.object({
  serviceType: z.enum(['FREIGHT_FORWARDING', 'CUSTOMS_BROKERAGE', 'SHIPPING', 'TRANSSHIPMENT']).optional(),
  shipmentMode: z.enum(['AIR', 'SEA', 'ROAD', 'RAIL']).optional(),
  origin: z.string().min(2, 'Origin is required').optional(),
  destination: z.string().min(2, 'Destination is required').optional(),
  rate: z.number().positive('Rate must be positive').optional(),
  description: z.string().min(10, 'Description is too short').optional(),
}).refine((data) => {
  // At least one field must be present
  return Object.keys(data).length > 0;
}, 'At least one field must be provided');

export const bookingValidationSchema = z.object({
  serviceId: z.string().uuid('Invalid service ID').optional(),
  cargoDetails: z.object({
    weight: z.number().positive('Weight must be positive').optional(),
    volume: z.number().positive('Volume must be positive').optional(),
    type: z.string().min(2, 'Cargo type is required').optional(),
    description: z.string().min(10, 'Description is too short').optional(),
    hazardous: z.boolean().optional(),
    specialInstructions: z.string().optional(),
  }).optional(),
}).refine((data) => {
  // At least one field must be present
  return Object.keys(data).length > 0;
}, 'At least one field must be provided');

// Input sanitization functions
export const sanitizeInput = (input: string): string => {
  return input
    .replace(/[<>]/g, '') // Remove < and > to prevent HTML injection
    .replace(/&/g, '&amp;')
    .replace(/"/g, '"')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;')
    .trim();
};

export const sanitizeObject = <T extends object>(obj: T): T => {
  const sanitized = { ...obj };
  
  for (const key in sanitized) {
    if (Object.prototype.hasOwnProperty.call(sanitized, key)) {
      const value = (sanitized as any)[key];
      
      if (typeof value === 'string') {
        (sanitized as any)[key] = sanitizeInput(value);
      } else if (typeof value === 'object' && value !== null) {
        (sanitized as any)[key] = sanitizeObject(value);
      }
    }
  }

  return sanitized;
};

// Validation helper functions
export const validateForm = async <T extends object>(
  schema: z.ZodSchema<T>,
  data: T
): Promise<{ success: boolean; errors?: string[] }> => {
  try {
    // First sanitize the input
    const sanitizedData = sanitizeObject(data);
    
    // Then validate against schema
    await schema.parseAsync(sanitizedData);
    return { success: true };
  } catch (error) {
    if (error instanceof z.ZodError) {
      return {
        success: false,
        errors: error.errors.map((e: z.ZodIssue) => e.message),
      };
    }
    return {
      success: false,
      errors: ['Validation failed'],
    };
  }
};

// Custom validators
export const validators = {
  isStrongPassword: (password: string): boolean => {
    return VALIDATION.PASSWORD_REGEX.test(password);
  },

  isValidEmail: (email: string): boolean => {
    return VALIDATION.EMAIL_REGEX.test(email);
  },

  isValidPhone: (phone: string): boolean => {
    return VALIDATION.PHONE_REGEX.test(phone);
  },

  isValidWalletAddress: (address: string): boolean => {
    return /^[A-Z0-9]{56}$/.test(address);
  },

  isValidAmount: (amount: string): boolean => {
    return /^\d+(\.\d{0,7})?$/.test(amount) && parseFloat(amount) > 0;
  },

  isValidURL: (url: string): boolean => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  },

  isValidJSON: (json: string): boolean => {
    try {
      JSON.parse(json);
      return true;
    } catch {
      return false;
    }
  },

  hasNoSQLInjection: (input: string): boolean => {
    const sqlInjectionPattern = /('|"|;|--|\/\*|\*\/|@@|@|\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|TABLE|OR|AND)\b)/i;
    return !sqlInjectionPattern.test(input);
  },

  hasNoXSS: (input: string): boolean => {
    const xssPattern = /<[^>]*>|javascript:|data:|vbscript:|on\w+=/i;
    return !xssPattern.test(input);
  },
};

// Export validation utilities
export const validation = {
  validateForm,
  sanitizeInput,
  sanitizeObject,
  validators,
};

export default validation;
