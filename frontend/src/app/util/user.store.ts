import { Injectable, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { jwtVerify } from 'jose';

export interface UserProfile {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  dataLimitMB: number | null;
  exp: number;
  iat: number;
}

@Injectable({ providedIn: 'root' })
export class UserStore {
  private router = inject(Router);

  // Core reactive state
  private readonly _user = signal<UserProfile | null>(null);
  readonly user = this._user.asReadonly();

  // Derived signals for templates
  readonly isAuthenticated = computed(() => this._user() !== null);
  readonly fullName = computed(() => {
    const u = this._user();
    return u ? `${u.firstName} ${u.lastName}`.trim() : '';
  });

  private cachedCryptoKey: CryptoKey | null = null;

  /**
   * Reads 'session_token', verifies its ES256 signature against your hex public key,
   * and populates the user signal. Returns true if valid.
   */
  async loadAndVerifySession(): Promise<boolean> {
    const token = this.getCookie('session_token');
    if (!token) {
      this._user.set(null);
      return false;
    }

    try {
      const publicKey = await this.getPublicKey();
      const { payload } = await jwtVerify(token, publicKey, {
        algorithms: ['ES256']
      });

      this._user.set({
        id: String(payload['id'] ?? ''),
        email: String(payload['email'] ?? ''),
        firstName: String(payload['firstName'] ?? ''),
        lastName: String(payload['lastName'] ?? ''),
        dataLimitMB: typeof payload['dataLimitMB'] === 'number' ? payload['dataLimitMB'] : null,
        exp: Number(payload['exp'] ?? 0),
        iat: Number(payload['iat'] ?? 0)
      });

      return true;
    } catch {
      this._user.set(null);
      return false;
    }
  }

  /**
   * Clears the session_token cookie, resets state, and redirects to /login.
   */
  logout(): void {
    document.cookie = 'session_token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
    this._user.set(null);
    this.router.navigate(['/login']);
  }

  private getCookie(name: string): string | null {
    if (typeof document === 'undefined') return null;
    const match = document.cookie.match(new RegExp('(^|;\\s*)' + name + '=([^;]*)'));
    return match ? decodeURIComponent(match[2]) : null;
  }

  private async getPublicKey(): Promise<CryptoKey> {
    if (this.cachedCryptoKey) return this.cachedCryptoKey;

    const hexKey = import.meta.env['NG_APP_JWT_PUBLIC_KEY'] ?? '';
    const cleanHex = hexKey.trim().replace(/^0x/i, '');
    const bytes = new Uint8Array(cleanHex.length / 2);
    for (let i = 0; i < cleanHex.length; i += 2) {
      bytes[i / 2] = parseInt(cleanHex.substring(i, i + 2), 16);
    }

    this.cachedCryptoKey = await crypto.subtle.importKey(
      'raw',
      bytes.buffer as ArrayBuffer,
      { name: 'ECDSA', namedCurve: 'P-256' },
      true,
      ['verify']
    );

    return this.cachedCryptoKey;
  }
}
