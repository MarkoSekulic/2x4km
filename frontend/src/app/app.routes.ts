
import { Routes } from '@angular/router';
import { Component } from '@angular/core';

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
  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full'
  },
  {
    path: 'dashboard',
    component: DashboardComponent
  }
];
