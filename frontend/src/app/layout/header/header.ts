import {
  Component,
  ElementRef,
  HostListener,
  OnInit,
  computed,
  inject,
  signal
} from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  LucideSearch,
  LucideBell,
  LucideChevronDown,
  LucideLogOut,
  LucideSettings,
  LucideShare2,
  LucideFilePlus,
  LucideRefreshCw
} from '@lucide/angular';
import { UserStore } from '../../util/user.store';
import { AppNotification, NotificationService } from '../../services/notification.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    RouterLink,
    LucideSearch,
    LucideBell,
    LucideChevronDown,
    LucideLogOut,
    LucideSettings,
    LucideShare2,
    LucideFilePlus,
    LucideRefreshCw
  ],
  templateUrl: './header.html',
  styleUrl: './header.css'
})
export class HeaderComponent implements OnInit {
  protected readonly userStore = inject(UserStore);
  private readonly notificationService = inject(NotificationService);
  private readonly elementRef = inject(ElementRef);

  profileDropdownOpen = signal(false);
  notificationsOpen = signal(false);
  isLoadingNotifications = signal(false);
  notifications = signal<AppNotification[]>([]);

  initials = computed(() => {
    const u = this.userStore.user();
    if (!u) return 'U';
    const first = u.firstName?.charAt(0) ?? '';
    const last = u.lastName?.charAt(0) ?? '';
    return (first + last).toUpperCase() || 'U';
  });

  ngOnInit(): void {
    this.loadNotifications();
  }

  loadNotifications(): void {
    this.isLoadingNotifications.set(true);
    this.notificationService.getNotifications().subscribe({
      next: (items) => {
        const sorted = [...items].sort(
          (a, b) => new Date(b.notifiedAt).getTime() - new Date(a.notifiedAt).getTime()
        );
        this.notifications.set(sorted);
        this.isLoadingNotifications.set(false);
      },
      error: () => {
        this.isLoadingNotifications.set(false);
      }
    });
  }

  toggleNotifications(): void {
    const willOpen = !this.notificationsOpen();
    this.notificationsOpen.set(willOpen);
    this.profileDropdownOpen.set(false);
    if (willOpen) {
      this.loadNotifications();
    }
  }

  toggleProfileDropdown(): void {
    this.profileDropdownOpen.update((open) => !open);
    this.notificationsOpen.set(false);
  }

  logout(): void {
    this.profileDropdownOpen.set(false);
    this.userStore.logout();
  }

  getTargetRoute(item: AppNotification): string[] {
    return item.fType === 'folder' ? ['/folders'] : ['/files'];
  }

  getIconBg(type: string): string {
    switch (type) {
      case 'shared':
        return '#eff6ff';
      case 'added':
        return '#f0fdf4';
      case 'updated':
        return '#fefce8';
      default:
        return '#f1f5f9';
    }
  }

  getIconColor(type: string): string {
    switch (type) {
      case 'shared':
        return '#2563eb';
      case 'added':
        return '#16a34a';
      case 'updated':
        return '#ca8a04';
      default:
        return '#475569';
    }
  }

  formatRelativeTime(isoString: string): string {
    if (!isoString) return '';
    const diffMs = Date.now() - new Date(isoString).getTime();
    const minutes = Math.floor(diffMs / 60000);
    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 7) return `${days}d ago`;
    return new Date(isoString).toLocaleDateString();
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.profileDropdownOpen.set(false);
      this.notificationsOpen.set(false);
    }
  }
}
