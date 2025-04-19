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

## TDD フロー

IMPORTANT! テスト駆動開発を実践し、以下のサイクルに従ってください:

1. **Red**: 失敗するテストを書く

   - 実装する機能の要件を明確に理解する
   - 期待される動作を検証するテストを書く
   - テストが失敗することを確認する

2. **Green**: テストが通るコードを実装する

   - 最小限の実装でテストを通す
   - 完璧なコードにこだわらない（リファクタリングは次のステップで）

3. **Refactor**: コードをクリーンアップする
   - テストが通ることを維持しながらコードを改善
   - コードの重複を取り除く
   - 命名を明確にする
   - パフォーマンスを向上させる

### テストファイル命名規則

- バックエンド: `{テスト対象ファイル名}_test.go`
- フロントエンド: `{テスト対象コンポーネント名}.test.tsx`

### モック作成ガイドライン

- バックエンド: `gomock`ライブラリを使用
- フロントエンド: `jest.mock()`または MSW を使用した API モック
- モックはテストディレクトリの`mocks`サブディレクトリに配置

## CI/CD パイプライン

プロジェクトは以下の CI/CD パイプラインを使用しています:

```
GitHub Push → テスト実行 → Lint/型チェック → ビルド → デプロイ
```

### GitHub Actions 設定

- `.github/workflows/ci.yml`: プルリクエスト時の検証
- `.github/workflows/deploy.yml`: メインブランチへのマージ時のデプロイ

### デプロイ環境

- 開発環境: `dev.guideforge.example.com`
- ステージング環境: `staging.guideforge.example.com`
- 本番環境: `guideforge.example.com`

### デプロイ手順

1. メインブランチにマージ
2. GitHub Actions による自動デプロイ
3. デプロイ後の検証

## デバッグテクニック

### フロントエンド

- React Developer Tools を使用してコンポーネント状態を確認
- Next.js の`useDebugValue`フックでカスタム値をデバッグ
- Chrome の Network タブで API リクエストを確認

デバッグコマンド:

```bash
# 詳細なログを有効にして開発サーバーを起動
DEBUG=* npm run dev
```

### バックエンド

- ログレベルを`debug`に設定してより詳細なログを確認
- `dlv`デバッガーを使用してコードをステップ実行

デバッグコマンド:

```bash
# デバッグモードでサーバー起動
LOG_LEVEL=debug go run cmd/api/main.go

# dlvでデバッグ
dlv debug cmd/api/main.go
```

### 一般的な問題と解決策

- JWT 認証エラー: トークンの有効期限とシークレットキーを確認
- 画像アップロード失敗: ストレージパーミッションとファイルサイズ制限を確認
- 遅い API レスポンス: SQL クエリの INDEX と N+1 問題を確認

## パフォーマンス最適化

### データベース最適化

- クエリには INDEX を適切に設定する
- 大量データ取得には`LIMIT`と`OFFSET`を使用する
- 複雑な集計クエリはビューを検討する

### フロントエンド最適化

- コンポーネントのメモ化（React.memo, useMemo, useCallback）
- 画像の最適化（next/image 使用）
- Code Splitting と遅延ロード

### キャッシュ戦略

- バックエンド: Redis を使用した結果キャッシュ
- フロントエンド: React Query のキャッシング
- CDN による静的アセットのキャッシュ

## API ドキュメント

API ドキュメントは Swagger で自動生成されています。

アクセス方法:

```
開発環境: http://localhost:8080/swagger/index.html
```

### 主要エンドポイント

- 認証: `/api/v1/auth/*`
- マニュアル: `/api/v1/manuals/*`
- ユーザー: `/api/v1/users/*`
- 画像アップロード: `/api/v1/upload/*`

### API レスポンス形式

すべての API レスポンスは以下の統一フォーマットに従います:

```json
{
  "success": true/false,
  "data": {...},
  "error": "エラーメッセージ"
}
```

## リリースプロセス

### バージョニング規則

セマンティックバージョニングを採用しています:

- メジャーバージョン: 互換性のない変更
- マイナーバージョン: 後方互換性のある機能追加
- パッチバージョン: バグ修正

### リリースチェックリスト

YOU MUST リリース前に以下の項目を確認してください:

1. すべてのテストが通過していること
2. セキュリティスキャンでの問題がないこと
3. パフォーマンステストで許容範囲内であること
4. リリースノートが作成されていること
5. データベースマイグレーションが確認されていること

### ホットフィックス手順

1. `hotfix/説明`ブランチを作成
2. 修正をコミット
3. テストを実行
4. 本番とメインブランチにマージ

## アーキテクチャ決定記録（ADR）

### ADR-001: フロントエンドフレームワークとして Next.js を採用

- **決定**: UI フレームワークとして Next.js（App Router）を使用
- **状況**: マニュアル管理システムに適した UI フレームワークが必要
- **選択した理由**: SEO 対応、SSR/SSG 対応、React Server Components によるパフォーマンス向上
- **検討した選択肢**: CRA, Nuxt.js, Angular
- **結果**: 高速なページロード、良好な DX、TypeScript との相性が確認された

### ADR-002: バックエンドフレームワークとして Echo を採用

- **決定**: Golang ウェブフレームワークとして Echo を使用
- **状況**: 高性能な REST API が必要
- **選択した理由**: ミドルウェアサポート、パフォーマンス、シンプルな API
- **検討した選択肢**: Gin, Fiber, net/http
- **結果**: 開発効率の向上と十分なパフォーマンスを実現

### ADR-003: 認証に JWT を採用

- **決定**: ユーザー認証に JWT を使用
- **状況**: ステートレスな認証メカニズムが必要
- **選択した理由**: スケーラビリティ、フロントエンドとバックエンドの分離容易性
- **検討した選択肢**: セッションベース認証、OAuth
- **結果**: 分散システムでの認証情報の共有が容易になった
