import CryptoJS from 'crypto-js';

// Encryption key should be stored securely and not in the code
// In production, this should be injected through environment variables
const ENCRYPTION_KEY = process.env['REACT_APP_ENCRYPTION_KEY'] || 'default-key-for-development';

/**
 * Encrypt sensitive data using AES encryption
 * @param data - Data to encrypt
 * @returns Encrypted data as string
 */
export const encrypt = (data: string): string => {
  try {
    return CryptoJS.AES.encrypt(data, ENCRYPTION_KEY).toString();
  } catch (error) {
    console.error('Encryption error:', error);
    throw new Error('Failed to encrypt data');
  }
};

/**
 * Decrypt encrypted data
 * @param encryptedData - Encrypted data to decrypt
 * @returns Decrypted data as string
 */
export const decrypt = (encryptedData: string): string => {
  try {
    const bytes = CryptoJS.AES.decrypt(encryptedData, ENCRYPTION_KEY);
    return bytes.toString(CryptoJS.enc.Utf8);
  } catch (error) {
    console.error('Decryption error:', error);
    throw new Error('Failed to decrypt data');
  }
};

/**
 * Hash data using SHA256
 * @param data - Data to hash
 * @returns Hashed data as string
 */
export const hash = (data: string): string => {
  try {
    return CryptoJS.SHA256(data).toString();
  } catch (error) {
    console.error('Hashing error:', error);
    throw new Error('Failed to hash data');
  }
};

/**
 * Generate a secure random string
 * @param length - Length of the random string
 * @returns Random string
 */
export const generateSecureRandomString = (length: number = 32): string => {
  const array = new Uint8Array(length);
  crypto.getRandomValues(array);
  return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
};

/**
 * Securely compare two strings in constant time
 * @param a - First string
 * @param b - Second string
 * @returns True if strings are equal
 */
export const secureCompare = (a: string, b: string): boolean => {
  if (a.length !== b.length) {
    return false;
  }

  let result = 0;
  for (let i = 0; i < a.length; i++) {
    result |= a.charCodeAt(i) ^ b.charCodeAt(i);
  }
  return result === 0;
};

/**
 * Generate a secure key derivation using PBKDF2
 * @param password - Password to derive key from
 * @param salt - Salt for key derivation
 * @returns Derived key as string
 */
export const deriveKey = (password: string, salt: string): string => {
  try {
    return CryptoJS.PBKDF2(password, salt, {
      keySize: 256 / 32,
      iterations: 10000
    }).toString();
  } catch (error) {
    console.error('Key derivation error:', error);
    throw new Error('Failed to derive key');
  }
};

/**
 * Encrypt an object
 * @param obj - Object to encrypt
 * @returns Encrypted object with same structure
 */
export const encryptObject = <T extends object>(obj: T): T => {
  const encryptedObj = { ...obj };
  
  const sensitiveFields = [
    'password',
    'secretKey',
    'walletKey',
    'privateKey',
    'token',
    'secret',
    'apiKey'
  ];

  for (const key in encryptedObj) {
    if (Object.prototype.hasOwnProperty.call(encryptedObj, key)) {
      const value = (encryptedObj as any)[key];
      
      if (sensitiveFields.includes(key) && typeof value === 'string') {
        (encryptedObj as any)[key] = encrypt(value);
      } else if (typeof value === 'object' && value !== null) {
        (encryptedObj as any)[key] = encryptObject(value);
      }
    }
  }

  return encryptedObj;
};

/**
 * Decrypt an object
 * @param obj - Object to decrypt
 * @returns Decrypted object with same structure
 */
export const decryptObject = <T extends object>(obj: T): T => {
  const decryptedObj = { ...obj };
  
  const sensitiveFields = [
    'password',
    'secretKey',
    'walletKey',
    'privateKey',
    'token',
    'secret',
    'apiKey'
  ];

  for (const key in decryptedObj) {
    if (Object.prototype.hasOwnProperty.call(decryptedObj, key)) {
      const value = (decryptedObj as any)[key];
      
      if (sensitiveFields.includes(key) && typeof value === 'string') {
        try {
          (decryptedObj as any)[key] = decrypt(value);
        } catch {
          // If decryption fails, the value might not be encrypted
          // Keep the original value
        }
      } else if (typeof value === 'object' && value !== null) {
        (decryptedObj as any)[key] = decryptObject(value);
      }
    }
  }

  return decryptedObj;
};

/**
 * Secure storage wrapper for sensitive data
 */
export const secureStorage = {
  set: (key: string, value: string): void => {
    const encryptedValue = encrypt(value);
    localStorage.setItem(key, encryptedValue);
  },

  get: (key: string): string | null => {
    const encryptedValue = localStorage.getItem(key);
    if (!encryptedValue) return null;
    try {
      return decrypt(encryptedValue);
    } catch {
      return null;
    }
  },

  remove: (key: string): void => {
    localStorage.removeItem(key);
  },

  clear: (): void => {
    localStorage.clear();
  }
};

export default {
  encrypt,
  decrypt,
  hash,
  generateSecureRandomString,
  secureCompare,
  deriveKey,
  encryptObject,
  decryptObject,
  secureStorage
};
