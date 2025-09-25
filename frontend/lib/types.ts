export interface Profile {
  id: string
  email: string
  display_name: string
  avatar_url?: string
  storage_used: number
  storage_limit: number
  created_at: string
  updated_at: string
}

export interface Folder {
  id: string
  name: string
  path: string
  parent_id?: string
  user_id: string
  created_at: string
  updated_at: string
}

export interface FileItem {
  id: string
  name: string
  original_name: string
  file_type: string
  mime_type: string
  file_size: number
  blob_url: string
  thumbnail_url?: string
  folder_id?: string
  user_id: string
  is_favorite: boolean
  ocr_content?: string
  video_duration?: number
  created_at: string
  updated_at: string
}

export interface FileShare {
  id: string
  file_id: string
  shared_by: string
  shared_with?: string
  is_public: boolean
  share_token: string
  expires_at?: string
  created_at: string
}

export type ViewMode = "grid" | "list"
export type SortBy = "name" | "date" | "size" | "type"
export type SortOrder = "asc" | "desc"
