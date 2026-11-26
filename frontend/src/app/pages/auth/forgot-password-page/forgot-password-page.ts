import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { AuthService } from '../auth'; // Adjust path as needed

@Component({
  selector: 'app-forgot-password',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <div class="auth-container">
      <h2>Reset Password</h2>

      @if (errorMessage()) {
        <div class="error-banner" role="alert">
          {{ errorMessage() }}
        </div>
      }

      @if (emailSent()) {
        <div class="success-box">
          <h3>Check your email</h3>
          <p>
            If an account exists for
            <strong>{{ forgotForm.controls.email.value }}</strong>,
            we sent a password reset link to your inbox.
          </p>
          <a routerLink="/login">Back to Sign In</a>
        </div>
      } @else {
        <form [formGroup]="forgotForm" (ngSubmit)="submitRequest()">
          <p class="info-text">
            Enter your email address and we'll send you a link to reset your password.
          </p>

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

          <button type="submit" [disabled]="forgotForm.invalid || isLoading()">
            {{ isLoading() ? 'Sending link...' : 'Send Reset Link' }}
          </button>

          <p class="footer-link">
            Remembered your password? <a routerLink="/login">Sign In</a>
          </p>
        </form>
      }
    </div>
  `,
  styles: [`
    .auth-container {
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
    .success-box {
      padding: 1rem;
      background-color: rgb(50, 50, 100);
      border: 1px solid #bbf7d0;
      border-radius: 6px;
    }
    button {
      width: 100%;
      padding: 0.7rem;
      margin-top: 0.5rem;
      cursor: pointer;
    }
    .footer-link {
      margin-top: 0.75rem;
      font-size: 0.9rem;
    }
  `]
})
export class ForgotPasswordPage {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);

  isLoading = signal(false);
  emailSent = signal(false);
  errorMessage = signal<string | null>(null);

  forgotForm = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]]
  });

  submitRequest(): void {
    if (this.forgotForm.invalid || this.isLoading()) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { email } = this.forgotForm.getRawValue();

    this.authService.requestPasswordReset(email).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.emailSent.set(true);
      },
      error: (err) => {
        this.isLoading.set(false);
        this.errorMessage.set(err?.error?.message || 'Unable to send reset email. Please try again.');
      }
    });
  }
}
