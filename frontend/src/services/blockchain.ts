import { Server, Networks, Asset, Keypair, TransactionBuilder } from '@stellar/stellar-sdk';
import { BLOCKCHAIN_NETWORK, BLOCKCHAIN_HORIZON_URL, TOKEN_CODE, TOKEN_ISSUER } from '../constants';

class BlockchainService {
  private server: Server;
  private network: string;
  private tokenAsset: Asset;

  constructor() {
    this.server = new Server(BLOCKCHAIN_HORIZON_URL);
    this.network = BLOCKCHAIN_NETWORK === 'TESTNET' ? Networks.TESTNET : Networks.PUBLIC;
    this.tokenAsset = new Asset(TOKEN_CODE, TOKEN_ISSUER);
  }

  async getBalance(address: string): Promise<string> {
    try {
      const account = await this.server.loadAccount(address);
      const balance = account.balances.find(
        (b: any) => b.asset_code === TOKEN_CODE && b.asset_issuer === TOKEN_ISSUER
      );
      return balance ? balance.balance : '0';
    } catch (error) {
      console.error('Error fetching balance:', error);
      throw new Error('Failed to fetch balance');
    }
  }

  async transfer(
    fromSecretKey: string,
    toAddress: string,
    amount: string,
    memo?: string
  ): Promise<string> {
    try {
      const sourceKeypair = Keypair.fromSecret(fromSecretKey);
      const sourceAccount = await this.server.loadAccount(sourceKeypair.publicKey());

      const transaction = new TransactionBuilder(sourceAccount, {
        fee: '100',
        networkPassphrase: this.network,
      })
        .addOperation({
          type: 'payment',
          destination: toAddress,
          asset: this.tokenAsset,
          amount: amount.toString(),
        })
        .setTimeout(30)
        .build();

      transaction.sign(sourceKeypair);

      const result = await this.server.submitTransaction(transaction);
      return result.hash;
    } catch (error) {
      console.error('Error transferring tokens:', error);
      throw new Error('Failed to transfer tokens');
    }
  }

  async createEscrow(
    fromSecretKey: string,
    escrowAddress: string,
    amount: string,
    duration: number
  ): Promise<string> {
    try {
      const sourceKeypair = Keypair.fromSecret(fromSecretKey);
      const sourceAccount = await this.server.loadAccount(sourceKeypair.publicKey());

      const transaction = new TransactionBuilder(sourceAccount, {
        fee: '100',
        networkPassphrase: this.network,
      })
        .addOperation({
          type: 'payment',
          destination: escrowAddress,
          asset: this.tokenAsset,
          amount: amount.toString(),
        })
        .setTimeout(duration)
        .build();

      transaction.sign(sourceKeypair);

      const result = await this.server.submitTransaction(transaction);
      return result.hash;
    } catch (error) {
      console.error('Error creating escrow:', error);
      throw new Error('Failed to create escrow');
    }
  }

  async releaseEscrow(
    escrowSecretKey: string,
    destinationAddress: string,
    amount: string
  ): Promise<string> {
    try {
      const escrowKeypair = Keypair.fromSecret(escrowSecretKey);
      const escrowAccount = await this.server.loadAccount(escrowKeypair.publicKey());

      const transaction = new TransactionBuilder(escrowAccount, {
        fee: '100',
        networkPassphrase: this.network,
      })
        .addOperation({
          type: 'payment',
          destination: destinationAddress,
          asset: this.tokenAsset,
          amount: amount.toString(),
        })
        .setTimeout(30)
        .build();

      transaction.sign(escrowKeypair);

      const result = await this.server.submitTransaction(transaction);
      return result.hash;
    } catch (error) {
      console.error('Error releasing escrow:', error);
      throw new Error('Failed to release escrow');
    }
  }

  async createAccount(destinationPublicKey: string): Promise<string> {
    try {
      const sourceKeypair = Keypair.fromSecret(TOKEN_ISSUER);
      const sourceAccount = await this.server.loadAccount(sourceKeypair.publicKey());

      const transaction = new TransactionBuilder(sourceAccount, {
        fee: '100',
        networkPassphrase: this.network,
      })
        .addOperation({
          type: 'createAccount',
          destination: destinationPublicKey,
          startingBalance: '1.5', // Minimum balance for a new account
        })
        .setTimeout(30)
        .build();

      transaction.sign(sourceKeypair);

      const result = await this.server.submitTransaction(transaction);
      return result.hash;
    } catch (error) {
      console.error('Error creating account:', error);
      throw new Error('Failed to create account');
    }
  }

  async trustAsset(secretKey: string): Promise<string> {
    try {
      const sourceKeypair = Keypair.fromSecret(secretKey);
      const sourceAccount = await this.server.loadAccount(sourceKeypair.publicKey());

      const transaction = new TransactionBuilder(sourceAccount, {
        fee: '100',
        networkPassphrase: this.network,
      })
        .addOperation({
          type: 'changeTrust',
          asset: this.tokenAsset,
          limit: '1000000000', // Maximum token amount the account can hold
        })
        .setTimeout(30)
        .build();

      transaction.sign(sourceKeypair);

      const result = await this.server.submitTransaction(transaction);
      return result.hash;
    } catch (error) {
      console.error('Error establishing trust line:', error);
      throw new Error('Failed to establish trust line');
    }
  }

  generateKeypair(): { publicKey: string; secretKey: string } {
    const keypair = Keypair.random();
    return {
      publicKey: keypair.publicKey(),
      secretKey: keypair.secret(),
    };
  }

  isValidPublicKey(publicKey: string): boolean {
    try {
      Keypair.fromPublicKey(publicKey);
      return true;
    } catch {
      return false;
    }
  }

  isValidSecretKey(secretKey: string): boolean {
    try {
      Keypair.fromSecret(secretKey);
      return true;
    } catch {
      return false;
    }
  }
}

export const blockchain = new BlockchainService();
export default blockchain;
