import { Component, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { AuthService } from '../auth'; // Adjust path as needed

@Component({
  selector: 'app-reset-password',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <div class="auth-container">
      <h2>Set New Password</h2>

      @if (errorMessage()) {
        <div class="error-banner" role="alert">
          {{ errorMessage() }}
        </div>
      }

      @if (!passForgotToken()) {
        <p>Invalid or missing password reset token.</p>
        <a routerLink="/forgot-password">Request a new reset link</a>
      } @else if (resetSuccess()) {
        <div class="success-box">
          <h3>Password Updated</h3>
          <p>Your password has been reset. Redirecting to sign in...</p>
          <a routerLink="/login">Go to Sign In now</a>
        </div>
      } @else {
        <form [formGroup]="resetForm" (ngSubmit)="submitReset()">
          <div class="form-field">
            <label for="email">Confirm Your Email</label>
            <input
              id="email"
              type="email"
              formControlName="email"
              autocomplete="email"
              placeholder="you@example.com"
            />
          </div>

          <div class="form-field">
            <label for="newPassword">New Password</label>
            <input
              id="newPassword"
              type="password"
              formControlName="newPassword"
              autocomplete="new-password"
            />
          </div>

          <button type="submit" [disabled]="resetForm.invalid || isLoading()">
            {{ isLoading() ? 'Updating password...' : 'Reset Password' }}
          </button>
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
  `]
})
export class ResetPasswordPage implements OnInit {
  private fb = inject(FormBuilder);
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private authService = inject(AuthService);

  passForgotToken = signal<string | null>(null);
  isLoading = signal(false);
  resetSuccess = signal(false);
  errorMessage = signal<string | null>(null);

  resetForm = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    newPassword: ['', [Validators.required, Validators.minLength(8)]]
  });

  ngOnInit(): void {
    const token = this.route.snapshot.paramMap.get('passForgotToken');
    this.passForgotToken.set(token);
  }

  submitReset(): void {
    const token = this.passForgotToken();
    if (!token || this.resetForm.invalid || this.isLoading()) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { email, newPassword } = this.resetForm.getRawValue();

    this.authService.resetPassword({
      passForgotToken: token,
      email,
      newPassword
    }).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.resetSuccess.set(true);
        setTimeout(() => this.router.navigate(['/login']), 2000);
      },
      error: (err) => {
        this.isLoading.set(false);
        this.errorMessage.set(err?.error?.message || 'Password reset failed. The link may have expired.');
      }
    });
  }
}
