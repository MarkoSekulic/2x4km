
import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { DatePipe } from '@angular/common';
import { forkJoin, of, switchMap } from 'rxjs';

import { FilesService } from '../../services/files.service';
import { Folder, CloudFile } from '../../models/files.models';

@Component({
  selector: 'app-files',
  imports: [DatePipe],
  templateUrl: './files.html',
  styleUrl: './files.css',
})
export class Files implements OnInit {
  private readonly filesService = inject(FilesService);
  private readonly cdr = inject(ChangeDetectorRef);

  showCreateFolder = false;
  newFolderName = '';
  isCreatingFolder = false;

  showRenameModal = false;
  renameItemType: 'file' | 'folder' | null = null;
  renameItemId = '';
  renameItemName = '';
  isRenaming = false;
  renameErrorMessage = '';

  deletingItemId: string | null = null;
  deleteErrorMessage = '';

  currentFolder: Folder | null = null;
  breadcrumbs: Folder[] = [];

  folders: Folder[] = [];
  files: CloudFile[] = [];

  errorMessage = '';
  isLoading = false;

  ngOnInit(): void {
    this.loadRootFolder();
  }

  // Loads the user's root folder.
  loadRootFolder(): void {
    this.isLoading = true;
    this.errorMessage = '';

    this.filesService.readRootFolder().pipe(
      switchMap(root => this.loadFolderContents(root).pipe(
        switchMap(result => of({ root, result }))
      ))
    ).subscribe({
      next: ({ root, result }) => {
        this.currentFolder = root;
        this.breadcrumbs = [root];

        this.folders = result.folders;
        this.files = result.files;

        this.isLoading = false;
        this.cdr.markForCheck();
      },
      error: error => {
        console.error('Error loading root folder:', error);
        this.errorMessage = 'Failed to load folder contents.';
        this.isLoading = false;
        this.cdr.markForCheck();
      }
    });
  }

  // Loads the contents of a specific folder.
  private loadFolderContents(folder: Folder) {
    const folders$ = folder.childFolders.length
      ? forkJoin(
          folder.childFolders.map(id => this.filesService.readFolder(id))
        )
      : of([] as Folder[]);

    const files$ = folder.files.length
      ? forkJoin(
          folder.files.map(id => this.filesService.readFile(id))
        )
      : of([] as CloudFile[]);

    return forkJoin({
      folders: folders$,
      files: files$
    });
  }

  openCreateFolder(): void {
  this.newFolderName = '';
  this.showCreateFolder = true;
}

closeCreateFolder(): void {
  if (this.isCreatingFolder) return;

  this.showCreateFolder = false;
  this.newFolderName = '';
}

openRenameFolder(folder: Folder): void {
  this.renameItemType = 'folder';
  this.renameItemId = folder.id;
  this.renameItemName = folder.name;
  this.renameErrorMessage = '';
  this.showRenameModal = true;
}

openRenameFile(file: CloudFile): void {
  this.renameItemType = 'file';
  this.renameItemId = file.id;
  this.renameItemName = file.name;
  this.renameErrorMessage = '';
  this.showRenameModal = true;
}

closeRenameModal(): void {
  if (this.isRenaming) return;

  this.showRenameModal = false;
  this.renameItemType = null;
  this.renameItemId = '';
  this.renameItemName = '';
  this.renameErrorMessage = '';
}

renameItem(): void {
  const newName = this.renameItemName.trim();

  if (
    !newName ||
    !this.renameItemId ||
    !this.renameItemType ||
    !this.currentFolder ||
    this.isRenaming
  ) {
    return;
  }

  const itemId = this.renameItemId;
  const itemType = this.renameItemType;
  const parentId = this.currentFolder.id;

  this.isRenaming = true;
  this.renameErrorMessage = '';

  const request$ = itemType === 'folder'
    ? this.filesService.updateFolder({
        id: itemId,
        newName,
        newParentId: parentId
      })
    : this.filesService.updateFile({
        id: itemId,
        newName,
        newParentId: parentId
      });

  request$.subscribe({
    next: () => {
      this.isRenaming = false;
      this.showRenameModal = false;
      this.renameItemType = null;
      this.renameItemId = '';
      this.renameItemName = '';

      this.cdr.markForCheck();
      this.refresh();
    },
    error: error => {
      console.error('Error renaming item:', error);

      this.renameErrorMessage = 'Failed to rename item.';
      this.isRenaming = false;

      this.cdr.markForCheck();
    }
  });
}

createFolder(): void {
  const name = this.newFolderName.trim();

  if (!name || !this.currentFolder || this.isCreatingFolder) {
    return;
  }

  this.isCreatingFolder = true;
  this.errorMessage = '';

  this.filesService.createFolder({
    name,
    parentId: this.currentFolder.id
  }).subscribe({
    next: () => {
      this.isCreatingFolder = false;
      this.showCreateFolder = false;
      this.newFolderName = '';

      this.cdr.markForCheck();
      this.refresh();
    },
    error: error => {
      console.error('Error creating folder:', error);

      this.errorMessage = 'Failed to create folder.';
      this.isCreatingFolder = false;

      this.cdr.markForCheck();
    }
  });
}

  deleteFolder(folder: Folder): void {
  if (this.deletingItemId) return;

  const confirmed = confirm(
    `Are you sure you want to delete the folder "${folder.name}"?`
  );

  if (!confirmed) return;

  this.deletingItemId = folder.id;
  this.deleteErrorMessage = '';

  this.filesService.deleteFolder(folder.id).subscribe({
    next: () => {
      this.deletingItemId = null;
      this.cdr.markForCheck();
      this.refresh();
    },
    error: error => {
      console.error('Error deleting folder:', error);

      this.deleteErrorMessage =
        'Failed to delete folder. Make sure it is empty.';

      this.deletingItemId = null;
      this.cdr.markForCheck();
    }
  });
}

deleteFile(file: CloudFile): void {
  if (this.deletingItemId) return;

  const confirmed = confirm(
    `Are you sure you want to delete the file "${file.name}"?`
  );

  if (!confirmed) return;

  this.deletingItemId = file.id;
  this.deleteErrorMessage = '';

  this.filesService.deleteFile(file.id).subscribe({
    next: () => {
      this.deletingItemId = null;
      this.cdr.markForCheck();
      this.refresh();
    },
    error: error => {
      console.error('Error deleting file:', error);

      this.deleteErrorMessage = 'Failed to delete file.';

      this.deletingItemId = null;
      this.cdr.markForCheck();
    }
  });
}

  // Opens a child folder.
  openFolder(folder: Folder): void {
    if (this.isLoading) return;

    const newBreadcrumbs = [...this.breadcrumbs, folder];

    this.loadFolder(folder, newBreadcrumbs);
  }

  // Navigates to a folder from the breadcrumb.
  navigateTo(index: number): void {
    if (this.isLoading) return;

    const folder = this.breadcrumbs[index];

    if (!folder) return;

    const newBreadcrumbs = this.breadcrumbs.slice(0, index + 1);

    this.loadFolder(folder, newBreadcrumbs);
  }

  // Loads a folder and updates the current navigation state.
  private loadFolder(folder: Folder, newBreadcrumbs: Folder[]): void {
    this.isLoading = true;
    this.errorMessage = '';

    this.loadFolderContents(folder).subscribe({
      next: result => {
        this.currentFolder = folder;
        this.breadcrumbs = newBreadcrumbs;

        this.folders = result.folders;
        this.files = result.files;

        this.isLoading = false;
        this.cdr.markForCheck();
      },
      error: error => {
        console.error('Error loading folder:', error);
        this.errorMessage = 'Failed to load folder contents.';
        this.isLoading = false;
        this.cdr.markForCheck();
      }
    });
  }

  // Refreshes the currently opened folder.
  refresh(): void {
  if (!this.currentFolder || this.isLoading) return;

  const folderId = this.currentFolder.id;
  const isRoot = this.breadcrumbs.length === 1;

  this.isLoading = true;
  this.errorMessage = '';

  const folder$ = isRoot
    ? this.filesService.readRootFolder()
    : this.filesService.readFolder(folderId);

  folder$.pipe(
    switchMap(folder =>
      this.loadFolderContents(folder).pipe(
        switchMap(result => of({ folder, result }))
      )
    )
  ).subscribe({
    next: ({ folder, result }) => {
      this.currentFolder = folder;

      this.breadcrumbs = this.breadcrumbs.map((item, index) =>
        index === this.breadcrumbs.length - 1 ? folder : item
      );

      this.folders = result.folders;
      this.files = result.files;

      this.isLoading = false;
      this.cdr.markForCheck();
    },
    error: error => {
      console.error('Error refreshing folder:', error);

      this.errorMessage = 'Failed to load folder contents.';
      this.isLoading = false;

      this.cdr.markForCheck();
    }
  });
}
}
