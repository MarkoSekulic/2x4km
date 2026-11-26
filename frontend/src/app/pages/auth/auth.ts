import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface TFATokenResponse {
  TFAToken: string;
}

export interface PassResetRequest {
  passForgotToken: string;
  email: string;
  newPassword: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private readonly apiUrl = import.meta.env['NG_APP_API_URL'] || '';

  // 1. RegisterPreVerify
  register(data: RegisterRequest): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/user/register`, data);
  }

  // 2. Register (Magic Link Verification)
  verifyEmail(verifyId: string): Observable<void> {
    return this.http.get<void>(`${this.apiUrl}/user/verify/${encodeURIComponent(verifyId)}`, {
      withCredentials: true // In case verify also sets session-token
    });
  }

  // 3. Login Step 1 -> Returns TFAToken
  login(data: LoginRequest): Observable<TFATokenResponse> {
    return this.http.post<TFATokenResponse>(`${this.apiUrl}/user/login`, data);
  }

  // 4. Login Step 2 -> Sends passkey, backend sets 'session-token' cookie
  verifyTfaLogin(tfaToken: string, passkey: string): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/user/TFALogin/${encodeURIComponent(tfaToken)}`,
      { TFAToken: tfaToken, passkey },
      { withCredentials: true } // Required for browser to accept Set-Cookie
    );
  }

  // 5. PrePasswordReset
  requestPasswordReset(email: string): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/user/prePassForgot`, { email });
  }

  // 6. PasswordReset
  resetPassword(data: PassResetRequest): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/user/passForgot/${encodeURIComponent(data.passForgotToken)}`,
      data
    );
  }
}
