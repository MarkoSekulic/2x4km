
export interface Folder {
  id: string;
  name: string;
  parent_id: string;
  root: string;
  child_folders: string[];
  files: string[];
  created_at: string;
  edited: string;
  ver: number;
}

export interface File {
  id: string;
  name: string;
  parent_id: string;
  root: string;
  created_at: string;
  edited: string;
  ver: number;
  osid: string;
}


export interface Folder {
  id: string;
  name: string;
  parentId: string;
  root: string;
  childFolders: string[];
  files: string[];
  createdAt: string;
  edited: string;
  ver: number;
}

export interface CloudFile {
  id: string;
  name: string;
  parentId: string;
  root: string;
  createdAt: string;
  edited: string;
  ver: number;
  osid: string;
}

export interface FolderCreateRequest {
  name: string;
  parentId: string;
}

export interface FolderUpdateRequest {
  id: string;
  newName: string;
  newParentId: string;
}

export interface FileCreateRequest {
  name: string;
  parentId: string;
}

export interface FileUpdateRequest {
  id: string;
  newName: string;
  newParentId: string;
}
