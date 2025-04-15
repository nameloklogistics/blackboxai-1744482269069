export const blockchain = {
  getBalance: jest.fn(),
  transfer: jest.fn(),
  createEscrow: jest.fn(),
  releaseEscrow: jest.fn(),
  createAccount: jest.fn(),
  trustAsset: jest.fn(),
  generateKeypair: jest.fn(),
  isValidPublicKey: jest.fn(),
  isValidSecretKey: jest.fn(),
};

export default blockchain;
