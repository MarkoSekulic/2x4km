
import { Component } from '@angular/core';

import {
  LucideSearch,
  LucideBell,
  LucideChevronDown
} from '@lucide/angular';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    LucideSearch,
    LucideBell,
    LucideChevronDown
  ],
  template: `
    <header class="header">

      <div class="search-container">
        <svg lucideSearch [size]="18"></svg>

        <input
          type="text"
          placeholder="Search files and folders..."
          aria-label="Search files and folders"
        />
      </div>

      <div class="header-actions">

        <button class="icon-button" aria-label="Notifications">
          <svg lucideBell [size]="19"></svg>
        </button>

        <div class="profile">
          <div class="avatar">U</div>

          <div class="profile-info">
            <span class="profile-name">User</span>
            <span class="profile-role">Personal account</span>
          </div>

          <svg lucideChevronDown [size]="16"></svg>
        </div>

      </div>

    </header>
  `,
  styleUrl: './header.css'
})
export class HeaderComponent {}
