# GuideForge プロジェクトガイド

このドキュメントは、GuideForge プロジェクトの開発に関する重要な情報を提供します。

## コマンド一覧

### 開発環境

```bash
# 開発環境の起動
docker-compose up -d

# フロントエンド開発サーバー起動
cd frontend
npm install
npm run dev

# バックエンド開発サーバー起動
cd backend
go mod download
go run cmd/api/main.go
```

### フロントエンド

```bash
# 依存関係のインストール
cd frontend
npm install

# 開発サーバー起動
npm run dev

# TypeScriptの型チェック
npm run typecheck

# リント実行
npm run lint

# ビルド
npm run build
```

### バックエンド

```bash
# 依存関係のダウンロード
cd backend
go mod download

# APIサーバー起動
go run cmd/api/main.go

# テスト実行
go test ./...

# 特定のテスト実行（例）
go test ./internal/api/manual_handler_test.go

# マイグレーション適用
go run cmd/migration/main.go up
```

## コーディングスタイル

### フロントエンド (Next.js/TypeScript)

IMPORTANT! 以下のスタイルガイドラインを遵守してください:

- App Router の構造に従う
- コンポーネントは機能ごとにディレクトリ分割する
- ニューモフィズムデザインの UI スタイルを一貫して使用する
- CSS Modules を使用し、グローバルスタイルは最小限にする
- TypeScript の型定義は厳密に行う（any 型は避ける）
- Common JS ではなく ES Modules 構文を使用する（import と export）
- データフェッチには適切なローディング状態とエラーハンドリングを含める
- 非同期処理は try-catch でエラーハンドリングする

### バックエンド (Golang/Echo)

IMPORTANT! 以下の規約を厳守してください:

- レイヤードアーキテクチャに従う（ハンドラー → サービス → リポジトリ）
- 関数名はキャメルケース、パッケージ名はスネークケースを使用
- エラーハンドリングは明示的に行い、ラップする（`fmt.Errorf("failed to get user: %w", err)`）
- ルーティングは cmd/api/main.go で一元管理する
- ログは internal/utils/logger.go を使用する
- 設定値は internal/config/から取得する
- トランザクション処理はサービスレイヤーで行う
- インターフェースを活用してモックしやすい設計にする

## プロジェクト構造

### コアファイル

- `frontend/src/app` - Next.js の App Router 構成
- `frontend/src/components` - 共通 UI コンポーネント
- `frontend/src/services` - API 通信
- `backend/internal/api` - API ハンドラー
- `backend/internal/services` - ビジネスロジック
- `backend/internal/repository` - データアクセス
- `backend/internal/models` - データモデル
- `database/init` - DB の初期化スクリプト

## データモデル

### 主要なモデル関連

マニュアルと手順の関係：

- `Manual`: タイトル、説明、カテゴリ、ユーザー ID、公開フラグを持つ
- `Step`: マニュアル ID、順序番号、タイトル、内容を持つ
- `Image`: ステップ ID、ファイルパス、ファイル名を持つ

これらのモデルは`backend/internal/models`ディレクトリで定義されています。

## ワークフロー

- 機能開発は新しいブランチで行う（命名規則: `feature/機能名`）
- コード変更後は必ずテストを実行する
- プルリクエスト前に lint とフォーマットを確認する
- マイグレーションファイルは連番で作成する

## 開発環境セットアップ

### 前提条件

- Docker と Docker Compose
- Node.js (v18 以上)
- Go (v1.21 以上)

YOU MUST 開発を始める前に次のセットアップを完了していることを確認してください:

1. リポジトリのクローン
2. `docker-compose up -d`でデータベース起動
3. バックエンドとフロントエンドの依存関係をインストール
4. フロントエンド開発サーバーとバックエンド開発サーバーを別々のターミナルで起動

## 注意点

- マニュアルデータのバリデーションは必ずバックエンドでも行う
- 画像アップロードは 5MB までに制限する
- フロントエンドの状態管理は React Query 推奨
- API リクエストのレスポンスは統一フォーマットを使用する
- データベースクエリはインデックスの使用を意識する

## テスト

- バックエンドはユニットテストとインテグレーションテストを実装
- フロントエンドはコンポーネントテストを主に実装
- テスト実行は CI パイプラインを使用

## トラブルシューティング

- PostgreSQL に接続できない場合は`docker-compose down -v && docker-compose up -d`を試す
- フロントエンドのビルドエラーは`npm run clean && npm install`で解決できることが多い
- バックエンドのホットリロードには`air`を使用できる
