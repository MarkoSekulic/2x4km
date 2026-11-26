
import { Routes } from '@angular/router';
import { Component } from '@angular/core';
import { LayoutComponent } from './layout/layout';
import { authGuard } from './util/auth-guard';

@Component({
  standalone: true,
  template: `
    <div>
      <h1>Dashboard</h1>
      <p>Welcome to your cloud storage.</p>
    </div>
  `
})
class DashboardComponent {}

export const routes: Routes = [
  
  { path: 'login', loadComponent: () => import("./pages/auth/login-page/login-page").then((m)=>m.LoginPage) },
  { path: 'register', loadComponent: () => import("./pages/auth/register-page/register-page").then((m)=>m.RegisterPage) },
  { path: 'verify/:verifyId', loadComponent: () => import("./pages/auth/verify-email-page/verify-email-page").then((m)=>m.VerifyEmailPage) },
  { path: 'forgot-password', loadComponent: () => import("./pages/auth/forgot-password-page/forgot-password-page").then((m)=>m.ForgotPasswordPage) },
  { path: 'reset-password/:passForgotToken', loadComponent: () => import("./pages/auth/reset-password-page/reset-password-page").then((m)=>m.ResetPasswordPage) },

  {
    path: '',
    canActivate: [authGuard],
    children: [
      {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full',
      },
      {
        path: '',
        component: LayoutComponent,
        children: [
          {
            path: 'dashboard',
            component: DashboardComponent
          }
        ]
      }
    ]
  },
];
