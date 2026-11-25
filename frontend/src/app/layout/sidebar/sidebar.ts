
import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

import {
  LucideCloud,
  LucideLayoutDashboard,
  LucideFiles,
  LucideFolder,
  LucideClock,
  LucideTrash2,
  LucideSettings
} from '@lucide/angular';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [
    RouterLink,
    RouterLinkActive,
    LucideCloud,
    LucideLayoutDashboard,
    LucideFiles,
    LucideFolder,
    LucideClock,
    LucideTrash2,
    LucideSettings
  ],
  template: `
    <aside class="sidebar">

      <div class="brand">
        <div class="brand-icon">
          <svg lucideCloud [size]="22"></svg>
        </div>
        <span>CloudBox</span>
      </div>

      <div class="navigation">

        <p class="section-label">WORKSPACE</p>

        <a routerLink="/dashboard"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideLayoutDashboard [size]="19"></svg>
          <span>Dashboard</span>
        </a>

        <a routerLink="/files"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideFiles [size]="19"></svg>
          <span>My files</span>
        </a>

        <a routerLink="/folders"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideFolder [size]="19"></svg>
          <span>Folders</span>
        </a>

        <a routerLink="/recent"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideClock [size]="19"></svg>
          <span>Recent</span>
        </a>

        <a routerLink="/trash"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideTrash2 [size]="19"></svg>
          <span>Trash</span>
        </a>

      </div>

      <div class="sidebar-bottom">
        <a routerLink="/settings"
           routerLinkActive="active"
           class="nav-item">
          <svg lucideSettings [size]="19"></svg>
          <span>Settings</span>
        </a>

        <div class="storage-card">
          <div class="storage-header">
            <span>Storage</span>
            <span>65%</span>
          </div>

          <div class="storage-progress">
            <div class="storage-progress-fill"></div>
          </div>

          <p>6.5 GB of 10 GB used</p>
        </div>
      </div>

    </aside>
  `,
  styleUrl: './sidebar.css'
})
export class SidebarComponent {}
