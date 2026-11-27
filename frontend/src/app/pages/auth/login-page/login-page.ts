import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../auth'; 

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <div class="login-container">
      <h2>{{ tfaToken() ? 'Two-Factor Authentication' : 'Sign In' }}</h2>

      @if (errorMessage()) {
        <div class="error-banner" role="alert">
          {{ errorMessage() }}
        </div>
      }

      <!-- STEP 1: Email & Password -->
      @if (!tfaToken()) {
        <form [formGroup]="credentialsForm" (ngSubmit)="submitCredentials()">
          <div class="form-field">
            <label for="email">Email</label>
            <input
              id="email"
              type="email"
              formControlName="email"
              autocomplete="email"
              placeholder="you@example.com"
            />
          </div>

          <div class="form-field">
            <label for="password">Password</label>
            <input
              id="password"
              type="password"
              formControlName="password"
              autocomplete="current-password"
            />
          </div>

          <div class="links">
            <a routerLink="/forgot-password">Forgot password?</a>
          </div>

          <button type="submit" [disabled]="credentialsForm.invalid || isLoading()">
            {{ isLoading() ? 'Sending code...' : 'Continue' }}
          </button>

          <p class="footer-link">
            Don't have an account? <a routerLink="/register">Register</a>
          </p>
        </form>
      }

      <!-- STEP 2: Email Passkey Verification -->
      @if (tfaToken()) {
        <form [formGroup]="tfaForm" (ngSubmit)="submitPasskey()">
          <p class="info-text">
            We sent a verification passkey to
            <strong>{{ credentialsForm.controls.email.value }}</strong>.
          </p>

          <div class="form-field">
            <label for="passkey">Email Passkey</label>
            <input
              id="passkey"
              type="text"
              formControlName="passkey"
              autocomplete="one-time-code"
              placeholder="Enter passkey"
              autofocus
            />
          </div>

          <button type="submit" [disabled]="tfaForm.invalid || isLoading()">
            {{ isLoading() ? 'Verifying...' : 'Verify & Sign In' }}
          </button>

          <button
            type="button"
            class="secondary-btn"
            [disabled]="isLoading()"
            (click)="resetToCredentials()"
          >
            Back to Login
          </button>
        </form>
      }
    </div>
  `,
  styles: [`
    .login-container {
      max-width: 400px;
      margin: 4rem auto;
      padding: 2rem;
      border: 1px solid #e2e8f0;
      border-radius: 8px;
    }
    .form-field {
      display: flex;
      flex-direction: column;
      margin-bottom: 1rem;
    }
    .form-field label {
      margin-bottom: 0.35rem;
      font-weight: 500;
    }
    .form-field input {
      padding: 0.6rem;
      border: 1px solid #cbd5e1;
      border-radius: 4px;
    }
    .error-banner {
      padding: 0.75rem;
      margin-bottom: 1rem;
      background-color: #fee2e2;
      color: #b91c1c;
      border-radius: 4px;
    }
    button {
      width: 100%;
      padding: 0.7rem;
      margin-top: 0.5rem;
      cursor: pointer;
    }
    .secondary-btn {
      background: transparent;
      border: 1px solid #cbd5e1;
    }
    .links, .footer-link {
      margin-top: 0.75rem;
      font-size: 0.9rem;
    }
  `]
})
export class LoginPage {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  // State signals
  tfaToken = signal<string | null>(null);
  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  // Step 1 Form (LoginRequest)
  credentialsForm = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required]]
  });

  // Step 2 Form (TFALoginRequest.passkey)
  tfaForm = this.fb.nonNullable.group({
    passkey: ['', [Validators.required]]
  });

  submitCredentials(): void {
    if (this.credentialsForm.invalid || this.isLoading()) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { email, password } = this.credentialsForm.getRawValue();

    this.authService.login({ email, password }).subscribe({
      next: (res) => {
        this.isLoading.set(false);
        if (res.TFAToken) {
          this.tfaToken.set(res.TFAToken);
        } else {
          this.errorMessage.set('Invalid server response: missing 2FA token.');
        }
      },
      error: (err) => {
        this.isLoading.set(false);
        this.errorMessage.set(err?.error?.message || 'Invalid email or password.');
      }
    });
  }

  submitPasskey(): void {
    const token = this.tfaToken();
    if (!token || this.tfaForm.invalid || this.isLoading()) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { passkey } = this.tfaForm.getRawValue();

    this.authService.verifyTfaLogin(token, passkey).subscribe({
      next: () => {
        this.isLoading.set(false);
        // Backend set the 'session-token' cookie; authGuard will now verify it
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        this.isLoading.set(false);
        this.errorMessage.set(err?.error?.message || 'Invalid or expired passkey.');
      }
    });
  }

  resetToCredentials(): void {
    this.tfaToken.set(null);
    this.tfaForm.reset();
    this.errorMessage.set(null);
  }
}
