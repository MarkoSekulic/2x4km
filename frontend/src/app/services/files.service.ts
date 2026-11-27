
import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';


import {
  Folder,
  CloudFile,
  FolderCreateRequest,
  FolderUpdateRequest,
  FileCreateRequest,
  FileUpdateRequest
} from '../models/files.models';

@Injectable({
  providedIn: 'root'
})
export class FilesService {
  private readonly http = inject(HttpClient);

  private readonly apiUrl = import.meta.env['NG_APP_API_URL'] || '';

  // Folders

  readFolder(id: string): Observable<Folder> {
    return this.http.post<Folder>(
      `${this.apiUrl}/folders/read`,
      { id },
      { withCredentials: true }
    );
  }

  createFolder(request: FolderCreateRequest): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/folders/create`,
      request,
      { withCredentials: true }
    );
  }

  updateFolder(request: FolderUpdateRequest): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/folders/update`,
      request,
      { withCredentials: true }
    );
  }

  deleteFolder(id: string): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/folders/delete`,
      { id },
      { withCredentials: true }
    );
  }

  readRootFolder(): Observable<Folder> {
  return this.http.get<Folder>(
    `${this.apiUrl}/folders/root`,
    { withCredentials: true }
  );
}

  // Files

  readFile(id: string): Observable<CloudFile> {
    return this.http.post<CloudFile>(
      `${this.apiUrl}/files/read`,
      { id },
      { withCredentials: true }
    );
  }

  createFile(request: FileCreateRequest): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/files/create`,
      request,
      { withCredentials: true }
    );
  }

  updateFile(request: FileUpdateRequest): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/files/update`,
      request,
      { withCredentials: true }
    );
  }

  deleteFile(id: string): Observable<void> {
    return this.http.post<void>(
      `${this.apiUrl}/files/delete`,
      { id },
      { withCredentials: true }
    );
  }
}
