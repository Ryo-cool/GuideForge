# GuideForge

業務効率化のためのマニュアル作成・管理ツール

## 概要

GuideForge は、業務マニュアルや手順書を簡単に作成・管理できる Web アプリケーションです。直感的なインターフェースを通じて、テキストと画像を組み合わせたステップバイステップのマニュアルを作成できます。

### 主な機能

- ユーザー認証（登録・ログイン）
- マニュアルの作成・編集・削除
- マニュアル内の手順（ステップ）管理
- 手順への画像添付
- マニュアルの検索・フィルタリング
- モダンなニューモフィズムデザインの UI

## 技術スタック

### フロントエンド

- Next.js 14 (App Router)
- TypeScript
- CSS Modules

### バックエンド

- Golang
- Echo Framework
- JWT 認証

### データベース

- PostgreSQL

### 開発環境

- Docker / Docker Compose

## プロジェクト構成

詳細なプロジェクト構成は [implementation-plan.md](./implementation-plan.md) を参照してください。

## 開発環境セットアップ

### 前提条件

- Docker と Docker Compose がインストール済みであること
- Node.js (v18 以上)
- Go (v1.21 以上)

### セットアップ手順

1. リポジトリをクローンします

```bash
git clone https://github.com/Ryo-cool/guideforge.git
cd guideforge
```

2. Docker Compose で開発環境を起動します

```bash
docker-compose up -d
```

3. フロントエンド開発サーバーを起動します

```bash
cd frontend
npm install
npm run dev
```

4. バックエンド開発サーバーを起動します（別ターミナルで）

```bash
cd backend
go mod download
go run cmd/api/main.go
```

5. ブラウザで http://localhost:3000 にアクセスしてアプリケーションを確認します

## 開発ワークフロー

1. [implementation-plan.md](./implementation-plan.md) のタスクリストを参照して実装を進めます
2. フロントエンドとバックエンドの実装を並行して進めることができます
3. `.cursor-rules.json` に定義されたコーディング規約に従って実装を行います

## ディレクトリ構造

```
guideforge/
├── frontend/                  # Next.js フロントエンド
│   └── ...
├── backend/                   # Golang バックエンド
│   └── ...
├── docker/                    # Docker関連ファイル
│   └── ...
├── docker-compose.yml         # 開発環境設定
├── .gitignore
└── README.md
```

## コントリビューション

1. 新機能の追加やバグ修正は、新しいブランチを作成してください
2. コードの変更を行ったら、既存のテストを実行してください
3. プルリクエストを作成する前に、コードがフォーマットされていることを確認してください
4. プルリクエストにはテストと適切なドキュメントを含めてください
