import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { jwtVerify } from 'jose';

function getCookie(name: string): string | null {
  if (typeof document === 'undefined') return null;
  const match = document.cookie.match(new RegExp('(^|;\\s*)' + name + '=([^;]*)'));
  return match ? decodeURIComponent(match[2]) : null;
}

// Converts a hex string (e.g. "04a1b2...") into a Uint8Array
function hexToBytes(hex: string): Uint8Array {
  const cleanHex = hex.trim().replace(/^0x/i, '');
  const bytes = new Uint8Array(cleanHex.length / 2);
  for (let i = 0; i < cleanHex.length; i += 2) {
    bytes[i / 2] = parseInt(cleanHex.substring(i, i + 2), 16);
  }
  return bytes;
}

// Imports an uncompressed EC public key (0x04 || X || Y) into Web Crypto
async function importUncompressedEcKey(hexKey: string, namedCurve: 'P-256' | 'P-384' | 'P-521' = 'P-256'): Promise<CryptoKey> {
  const rawKeyBytes = hexToBytes(hexKey);
  return crypto.subtle.importKey(
    'raw',
    rawKeyBytes.buffer as ArrayBuffer,
    { name: 'ECDSA', namedCurve },
    true,
    ['verify']
  );
}

export const authGuard: CanActivateFn = async () => {
  const router = inject(Router);
  const token = getCookie('session_token');

  console.log('[AuthGuard] document.cookie:', document.cookie);
  console.log('[AuthGuard] extracted token:', token);

  if (!token) {
    console.warn('[AuthGuard] No session_token cookie found -> redirecting to /login');
    return router.createUrlTree(['/login']);
  }

  try {
    const hexKey = import.meta.env['NG_APP_JWT_PUBLIC_KEY'];
    console.log('[AuthGuard] hexKey loaded:', hexKey ? `${hexKey.slice(0, 10)}... (len: ${hexKey.length})` : 'UNDEFINED');

    const publicKey = await importUncompressedEcKey(hexKey, 'P-256');
    const verified = await jwtVerify(token, publicKey, {
      algorithms: ['ES256']
    });

    console.log('[AuthGuard] JWT verified successfully:', verified.payload);
    return true;
  } catch (err) {
    console.error('[AuthGuard] Verification failed:', err);
    return router.createUrlTree(['/login']);
  }
};
