// ユーザー関連の型定義
export interface User {
  id: number;
  username: string;
  email: string;
  profileImage?: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserLogin {
  email: string;
  password: string;
}

export interface UserRegister {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

// マニュアル関連の型定義
export interface Manual {
  id: number;
  title: string;
  description?: string;
  category?: string;
  userId: number;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
  steps?: Step[];
}

export interface ManualCreate {
  title: string;
  description?: string;
  category?: string;
  isPublic?: boolean;
}

export interface ManualUpdate {
  title?: string;
  description?: string;
  category?: string;
  isPublic?: boolean;
}

// 手順関連の型定義
export interface Step {
  id: number;
  manualId: number;
  orderNumber: number;
  title: string;
  content?: string;
  createdAt: string;
  updatedAt: string;
  images?: Image[];
}

export interface StepCreate {
  title: string;
  content?: string;
  orderNumber?: number;
}

export interface StepUpdate {
  title?: string;
  content?: string;
  orderNumber?: number;
}

export interface StepOrderUpdate {
  id: number;
  orderNumber: number;
}

// 画像関連の型定義
export interface Image {
  id: number;
  stepId: number;
  filePath: string;
  fileName: string;
  fileSize: number;
  mimeType: string;
  createdAt: string;
}

// APIレスポンス関連の型定義
export interface ApiResponse<T> {
  status: 'success' | 'error';
  data?: T;
  message?: string;
  errors?: Record<string, string[]>;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// 検索・フィルタリング関連の型定義
export interface ManualSearchParams {
  search?: string;
  category?: string;
  isPublic?: boolean;
  page?: number;
  limit?: number;
  sortBy?: 'title' | 'createdAt' | 'updatedAt';
  sortOrder?: 'asc' | 'desc';
} 