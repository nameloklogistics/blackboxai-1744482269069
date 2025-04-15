import '@testing-library/jest-dom';

declare global {
  namespace jest {
    interface Matchers<R> {
      toBeWithinRange(floor: number, ceiling: number): R;
    }
  }

  const jest: typeof import('@jest/globals')['jest'];
  const expect: typeof import('@jest/globals')['expect'];
  const describe: typeof import('@jest/globals')['describe'];
  const it: typeof import('@jest/globals')['it'];
  const test: typeof import('@jest/globals')['test'];
  const beforeAll: typeof import('@jest/globals')['beforeAll'];
  const beforeEach: typeof import('@jest/globals')['beforeEach'];
  const afterAll: typeof import('@jest/globals')['afterAll'];
  const afterEach: typeof import('@jest/globals')['afterEach'];

  interface Window {
    matchMedia: (query: string) => {
      matches: boolean;
      media: string;
      onchange: null;
      addListener: jest.Mock;
      removeListener: jest.Mock;
      addEventListener: jest.Mock;
      removeEventListener: jest.Mock;
      dispatchEvent: jest.Mock;
    };
    IntersectionObserver: jest.Mock;
    ResizeObserver: jest.Mock;
    scrollTo: jest.Mock;
    localStorage: {
      getItem: jest.Mock;
      setItem: jest.Mock;
      removeItem: jest.Mock;
      clear: jest.Mock;
    };
    sessionStorage: {
      getItem: jest.Mock;
      setItem: jest.Mock;
      removeItem: jest.Mock;
      clear: jest.Mock;
    };
    URL: {
      createObjectURL: jest.Mock;
      revokeObjectURL: jest.Mock;
    } & typeof URL;
    MutationObserver: jest.Mock;
    crypto: {
      getRandomValues: (arr: any) => any;
      subtle: {
        digest: jest.Mock;
      };
    };
  }

  interface TextEncoderStream {
    readable: ReadableStream;
    writable: WritableStream;
    encoding: string;
  }

  interface TextDecoderStream {
    readable: ReadableStream;
    writable: WritableStream;
    encoding: string;
    fatal: boolean;
    ignoreBOM: boolean;
  }

  interface AllowSharedBufferSource {
    readonly [Symbol.toStringTag]: string;
    slice(begin?: number, end?: number): AllowSharedBufferSource;
  }

  interface TextDecodeOptions {
    stream?: boolean;
  }

  interface Console {
    log: jest.Mock;
    error: jest.Mock;
    warn: jest.Mock;
    info: jest.Mock;
    debug: jest.Mock;
  }
}

declare module '@stellar/stellar-sdk' {
  export class Server {
    constructor(url: string);
    loadAccount: jest.Mock;
    submitTransaction: jest.Mock;
  }

  export const Networks: {
    TESTNET: string;
    PUBLIC: string;
  };

  export class Asset {
    constructor(code: string, issuer: string);
  }

  export const Keypair: {
    random: jest.Mock;
    fromSecret: jest.Mock;
  };

  export class TransactionBuilder {
    constructor(account: any, options: any);
    addOperation: jest.Mock;
    setTimeout: jest.Mock;
    build: jest.Mock;
    sign: jest.Mock;
  }

  export const Operation: {
    payment: jest.Mock;
    createAccount: jest.Mock;
  };
}

declare module 'jest-canvas-mock';
declare module 'whatwg-fetch';

export {};
