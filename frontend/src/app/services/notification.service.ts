import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';

export interface AppNotification {
  userId: string;
  notifiedAt: string; // ISO 8601 timestamp string from google.protobuf.Timestamp
  type: 'shared' | 'added' | 'updated' | string;
  byUserId: string;
  byUserName: string;
  fType: 'file' | 'folder' | string;
  fId: string;
  fName: string;
  addedFId?: string;
  addedFName?: string;
}

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private http = inject(HttpClient);
  private readonly apiUrl = import.meta.env['NG_APP_API_URL'] || '';

  getNotifications(): Observable<AppNotification[]> {
    return this.http
      .get<{ notifications?: any[] }>(`${this.apiUrl}/notification/me`, {
        withCredentials: true // Sends session_token cookie
      })
      .pipe(
        map((res) =>
          (res.notifications ?? []).map((n) => ({
            userId: n.userId ?? n.user_id ?? '',
            notifiedAt: n.notifiedAt ?? n.notified_at ?? '',
            type: n.type ?? '',
            byUserId: n.byUserId ?? n.by_user_id ?? '',
            byUserName: n.byUserName ?? n.by_user_name ?? 'Someone',
            fType: n.fType ?? n.f_type ?? 'file',
            fId: n.fId ?? n.f_id ?? '',
            fName: n.fName ?? n.f_name ?? 'Untitled',
            addedFId: n.addedFId ?? n.added_f_id ?? '',
            addedFName: n.addedFName ?? n.added_f_name ?? ''
          }))
        )
      );
  }
}
