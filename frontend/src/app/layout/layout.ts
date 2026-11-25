
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarComponent } from './sidebar/sidebar';
import { HeaderComponent } from './header/header';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [
    RouterOutlet,
    SidebarComponent,
    HeaderComponent
  ],
  template: `
    <div class="app-layout">

      <app-sidebar></app-sidebar>

      <div class="main-container">
        <app-header></app-header>

        <main class="page-content">
          <router-outlet></router-outlet>
        </main>
      </div>

    </div>
  `,
  styleUrl: './layout.css'
})
export class LayoutComponent {}
