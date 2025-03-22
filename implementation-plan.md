# GuideForge - マニュアル作成ツール実装計画

## プロジェクト概要

GuideForge は、業務効率化のためのマニュアル作成・管理ツールです。初期 MVP として必要最低限の機能を実装し、ユーザーフィードバックを得ながら機能拡張を図ります。

### 技術スタック

- **フロントエンド**: Next.js（TypeScript）、ニューモフィズムデザイン
- **バックエンド**: Golang (echo フレームワーク)、REST API
- **データベース**: PostgreSQL

## ディレクトリ構成

```
guideforge/
├── frontend/                  # Next.js フロントエンド
│   ├── public/                # 静的ファイル
│   ├── src/
│   │   ├── app/              # App Router構成
│   │   ├── components/       # 共通コンポーネント
│   │   │   ├── ui/           # UIコンポーネント
│   │   │   ├── layout/       # レイアウト関連コンポーネント
│   │   │   └── forms/        # フォーム関連コンポーネント
│   │   ├── hooks/            # カスタムフック
│   │   ├── lib/              # ユーティリティ関数
│   │   ├── services/         # API呼び出し関連
│   │   ├── types/            # TypeScript型定義
│   │   └── styles/           # グローバルスタイル
│   ├── .env.local            # 環境変数
│   ├── package.json
│   └── tsconfig.json
│
├── backend/                   # Golang バックエンド
│   ├── cmd/
│   │   └── api/              # APIエントリーポイント
│   ├── internal/
│   │   ├── api/              # APIハンドラー
│   │   ├── auth/             # 認証関連
│   │   ├── config/           # 設定
│   │   ├── middleware/       # ミドルウェア
│   │   ├── models/           # データモデル
│   │   ├── repository/       # データアクセス層
│   │   ├── services/         # ビジネスロジック
│   │   └── utils/            # ユーティリティ
│   ├── pkg/                   # 共有パッケージ
│   ├── migrations/            # データベースマイグレーション
│   ├── go.mod
│   └── go.sum
│
├── database/                  # データベース関連ファイル
│   └── init/                 # 初期化スクリプト
│
├── .cursor/                   # Cursorプロジェクトルール
│   └── rules/                # 各種開発規約
│
├── docker-compose.yml         # 開発環境設定
├── .gitignore
└── README.md
```

## データベース設計

### ユーザーテーブル

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  profile_image VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### マニュアルテーブル

```sql
CREATE TABLE manuals (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  is_public BOOLEAN DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 手順テーブル

```sql
CREATE TABLE steps (
  id SERIAL PRIMARY KEY,
  manual_id INTEGER NOT NULL REFERENCES manuals(id) ON DELETE CASCADE,
  order_number INTEGER NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 画像テーブル

```sql
CREATE TABLE images (
  id SERIAL PRIMARY KEY,
  step_id INTEGER REFERENCES steps(id) ON DELETE CASCADE,
  file_path VARCHAR(255) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  file_size INTEGER NOT NULL,
  mime_type VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## API エンドポイント設計

### ユーザー管理

- `POST /api/users` - ユーザー登録
- `POST /api/login` - ログイン
- `GET /api/users/me` - 自分のプロフィール取得
- `PUT /api/users/me` - プロフィール更新
- `POST /api/password/reset` - パスワードリセット要求
- `PUT /api/password/reset` - パスワード更新

### マニュアル管理

- `GET /api/manuals` - マニュアル一覧取得
- `POST /api/manuals` - マニュアル作成
- `GET /api/manuals/{id}` - マニュアル詳細取得
- `PUT /api/manuals/{id}` - マニュアル更新
- `DELETE /api/manuals/{id}` - マニュアル削除

### 手順管理

- `GET /api/manuals/{id}/steps` - 特定マニュアルの手順一覧取得
- `POST /api/manuals/{id}/steps` - 手順追加
- `PUT /api/steps/{id}` - 手順更新
- `DELETE /api/steps/{id}` - 手順削除
- `PUT /api/manuals/{id}/steps/order` - 手順順序更新

### 画像管理

- `POST /api/steps/{id}/images` - 画像アップロード
- `DELETE /api/images/{id}` - 画像削除

## 実装タスクチェックリスト

### 環境構築

- [ ] プロジェクトディレクトリ構成作成
- [ ] フロントエンド環境構築 (Next.js)
- [ ] バックエンド環境構築 (Golang/Echo)
- [ ] Docker 設定
- [ ] データベース設定 (PostgreSQL)

### バックエンド実装

- [ ] データベース接続設定
- [ ] モデル定義
- [ ] マイグレーション作成
- [ ] リポジトリ層実装
- [ ] サービス層実装
- [ ] 認証・認可機能実装
  - [ ] JWT 実装
  - [ ] ミドルウェア設定
- [ ] API エンドポイント実装
  - [ ] ユーザー管理 API
  - [ ] マニュアル管理 API
  - [ ] 手順管理 API
  - [ ] 画像アップロード API
- [ ] エラーハンドリング
- [ ] ロギング設定
- [ ] テスト実装

### フロントエンド実装

- [ ] プロジェクト設定
- [ ] ルーティング設定
- [ ] コンポーネント設計・実装
  - [ ] ニューモフィズム UI コンポーネント
  - [ ] 共通レイアウト
  - [ ] フォームコンポーネント
- [ ] 認証関連実装
  - [ ] ログイン画面
  - [ ] ユーザー登録画面
  - [ ] 認証フック
- [ ] ページ実装
  - [ ] ダッシュボード
  - [ ] マニュアル一覧
  - [ ] マニュアル作成/編集
  - [ ] マニュアル詳細表示
  - [ ] プロフィール管理
- [ ] API クライアント実装
- [ ] 状態管理
- [ ] エラーハンドリング
- [ ] レスポンシブデザイン対応
- [ ] テスト実装

### インテグレーション

- [ ] フロントエンド・バックエンド連携テスト
- [ ] エンドツーエンドテスト
- [ ] パフォーマンス最適化

### デプロイ・運用

- [ ] CI/CD 設定
- [ ] 本番環境構築
- [ ] モニタリング設定
- [ ] バックアップ戦略
