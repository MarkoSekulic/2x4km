import { Component, OnInit, inject, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { AuthService } from '../auth'; // Adjust path as needed

@Component({
  selector: 'app-verify-email',
  standalone: true,
  imports: [RouterLink],
  template: `
    <div class="verify-container">
      @switch (status()) {
        @case ('verifying') {
          <h2>Verifying your email...</h2>
          <p>Please wait while we activate your account.</p>
        }
        @case ('success') {
          <h2>Email Verified</h2>
          <p>Your account is now active. Redirecting to sign in...</p>
          <a routerLink="/login">Go to Sign In now</a>
        }
        @case ('error') {
          <h2>Verification Failed</h2>
          <p class="error-text">{{ errorMessage() }}</p>
          <a routerLink="/register">Back to Register</a>
        }
      }
    </div>
  `,
  styles: [`
    .verify-container {
      max-width: 400px;
      margin: 4rem auto;
      padding: 2rem;
      text-align: center;
      border: 1px solid #e2e8f0;
      border-radius: 8px;
    }
    .error-text {
      color: #b91c1c;
      margin-bottom: 1rem;
    }
  `]
})
export class VerifyEmailPage implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private authService = inject(AuthService);

  status = signal<'verifying' | 'success' | 'error'>('verifying');
  errorMessage = signal<string>('This verification link is invalid or has expired.');

  ngOnInit(): void {
    const verifyId = this.route.snapshot.paramMap.get('verifyId');

    if (!verifyId) {
      this.status.set('error');
      return;
    }

    this.authService.verifyEmail(verifyId).subscribe({
      next: () => {
        this.status.set('success');
        // Change to ['/dashboard'] if your backend sets session_token on /user/verify
        setTimeout(() => this.router.navigate(['/login']), 2000);
      },
      error: (err) => {
        this.status.set('error');
        if (err?.error?.message) {
          this.errorMessage.set(err.error.message);
        }
      }
    });
  }
}
