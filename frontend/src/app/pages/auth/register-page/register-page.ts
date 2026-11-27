import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { AuthService } from '../auth'; // Adjust path as needed

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <div class="auth-container">
      <h2>Create an Account</h2>

      @if (errorMessage()) {
        <div class="error-banner" role="alert">
          {{ errorMessage() }}
        </div>
      }

      @if (emailSent()) {
        <div class="success-box">
          <h3>Check your email</h3>
          <p>
            We sent a verification link to
            <strong>{{ registerForm.controls.email.value }}</strong>.
            Click the link in the email to activate your account.
          </p>
          <a routerLink="/login">Back to Sign In</a>
        </div>
      } @else {
        <form [formGroup]="registerForm" (ngSubmit)="submitRegister()">
          <div class="form-field">
            <label for="firstName">First Name</label>
            <input id="firstName" type="text" formControlName="firstName" autocomplete="given-name" />
          </div>

          <div class="form-field">
            <label for="lastName">Last Name</label>
            <input id="lastName" type="text" formControlName="lastName" autocomplete="family-name" />
          </div>

          <div class="form-field">
            <label for="email">Email</label>
            <input id="email" type="email" formControlName="email" autocomplete="email" />
          </div>

          <div class="form-field">
            <label for="password">Password</label>
            <input id="password" type="password" formControlName="password" autocomplete="new-password" />
          </div>

          <button type="submit" [disabled]="registerForm.invalid || isLoading()">
            {{ isLoading() ? 'Creating account...' : 'Register' }}
          </button>

          <p class="footer-link">
            Already have an account? <a routerLink="/login">Sign In</a>
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
export class RegisterPage {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);

  isLoading = signal(false);
  emailSent = signal(false);
  errorMessage = signal<string | null>(null);

  registerForm = this.fb.nonNullable.group({
    firstName: ['', [Validators.required]],
    lastName: ['', [Validators.required]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(8)]]
  });

  submitRegister(): void {
    if (this.registerForm.invalid || this.isLoading()) return;

    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.authService.register(this.registerForm.getRawValue()).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.emailSent.set(true);
      },
      error: (err) => {
        this.isLoading.set(false);
        this.errorMessage.set(err?.error?.message || 'Registration failed. Please try again.');
      }
    });
  }
}
