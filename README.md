# プライベートブロックチェーン社内データ管理システム

企業向けプライベートブロックチェーンを利用した安全な社内データ管理システムです。データの改ざん防止、透明性の確保、細かなアクセス制御を実現し、企業の重要なデータを安全に管理します。

## 特徴

- **プライベートブロックチェーン**: 全てのデータ操作がブロックチェーンに記録され、改ざんを防止
- **役割ベースアクセス制御 (RBAC)**: 管理者、マネージャー、従業員、ゲストの4段階の権限管理
- **データ暗号化**: AES-256による保存時暗号化とTLS通信
- **監査ログ**: 全ての操作が詳細に記録され、コンプライアンス要件に対応
- **バージョン管理**: ドキュメントの変更履歴を自動追跡
- **細かな権限設定**: ファイル・フォルダ単位での詳細なアクセス制御

## 技術スタック

### バックエンド
- **Go 1.21+** - メインアプリケーション言語
- **Gin** - HTTPウェブフレームワーク
- **GORM** - ORMライブラリ
- **PostgreSQL** - メインデータベース
- **Redis** - キャッシュ・セッション管理
- **JWT** - 認証トークン

### インフラ
- **Docker & Docker Compose** - コンテナ化
- **Nginx** - リバースプロキシ・ロードバランサー

## プロジェクト構造

```
├── cmd/                    # アプリケーションエントリーポイント
│   └── server/            # メインサーバー
├── internal/              # 内部パッケージ
│   ├── api/              # API関連
│   │   ├── handlers/     # HTTPハンドラー
│   │   ├── middleware/   # ミドルウェア
│   │   └── routes/       # ルーティング
│   ├── blockchain/       # ブロックチェーン実装
│   ├── config/           # 設定管理
│   ├── database/         # データベース関連
│   │   └── models/       # データモデル
│   ├── security/         # セキュリティ機能
│   │   ├── auth/         # 認証
│   │   ├── crypto/       # 暗号化
│   │   └── rbac/         # アクセス制御
│   └── services/         # ビジネスロジック
├── deployments/          # デプロイメント設定
├── docs/                 # ドキュメント
└── tests/                # テスト
```

## クイックスタート

### 前提条件

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (Dockerを使用しない場合)

### 1. リポジトリのクローン

```bash
git clone https://github.com/nshmdayo/in-house-datamanagement-system-sample.git
cd in-house-datamanagement-system-sample
```

### 2. 開発環境のセットアップ

```bash
make setup-dev
```

### 3. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集して設定を調整
```

### 4. 依存関係のインストール

```bash
make deps
```

### 5. Docker Composeでの起動

```bash
make docker-compose-up
```

### 6. ローカル開発での起動

```bash
# データベースが起動していることを確認
make dev  # ホットリロード付きで開発サーバー起動
```

## API エンドポイント

### 認証
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークン更新
- `POST /api/v1/auth/logout` - ログアウト
- `GET /api/v1/auth/profile` - プロファイル取得

### ユーザー管理
- `GET /api/v1/users` - ユーザー一覧
- `POST /api/v1/users` - ユーザー作成 (管理者のみ)
- `GET /api/v1/users/:id` - ユーザー詳細
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除 (管理者のみ)

### ドキュメント管理
- `GET /api/v1/documents` - ドキュメント一覧
- `POST /api/v1/documents` - ドキュメント作成
- `GET /api/v1/documents/:id` - ドキュメント詳細
- `PUT /api/v1/documents/:id` - ドキュメント更新
- `DELETE /api/v1/documents/:id` - ドキュメント削除

### ブロックチェーン
- `GET /api/v1/blockchain/blocks` - ブロック一覧
- `GET /api/v1/blockchain/transactions/:id` - トランザクション詳細
- `POST /api/v1/blockchain/verify` - データ整合性検証

### 監査ログ
- `GET /api/v1/audit/logs` - 監査ログ一覧 (管理者・マネージャーのみ)
- `GET /api/v1/audit/statistics` - 統計情報

## 開発コマンド

```bash
# 開発サーバー起動
make dev

# ビルド
make build

# テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# リント
make lint

# フォーマット
make fmt

# セキュリティスキャン
make security

# Docker イメージビルド
make docker-build

# データベースマイグレーション
make db-migrate

# データベースシード
make db-seed
```

## セキュリティ機能

### 認証・認可
- JWT ベースの認証
- 役割ベースアクセス制御 (RBAC)
- アカウントロック機能 (ログイン失敗時)
- セッション管理

### データ保護
- AES-256 による保存時暗号化
- TLS 1.3 による転送時暗号化
- bcrypt によるパスワードハッシュ化

### 監査・ログ
- 全操作の詳細ログ
- セキュリティイベントの追跡
- IP アドレス・User Agent の記録

## ブロックチェーン機能

### 特徴
- プライベートブロックチェーン実装
- Proof of Work コンセンサス
- Merkle Tree による効率的な検証
- データ整合性の自動検証

### 記録される操作
- ドキュメントの作成・更新・削除
- アクセス権限の変更
- ユーザー操作履歴

## テスト

```bash
# 全テスト実行
make test

# カバレッジレポート生成
make test-coverage

# ベンチマークテスト
make bench

# ロードテスト
make load-test
```

## デプロイメント

### Docker Compose (推奨)

```bash
# 本番環境での起動
make docker-compose-up

# ログ確認
make docker-compose-logs

# 停止
make docker-compose-down
```

### Kubernetes

```bash
# Kubernetes マニフェストを適用
kubectl apply -f deployments/k8s/
```

## モニタリング

### ヘルスチェック

```bash
curl http://localhost:8080/health
```

### メトリクス

アプリケーションメトリクスは `/metrics` エンドポイントで公開されています。

## セキュリティ考慮事項

- 定期的な依存関係の更新
- セキュリティパッチの適用
- 監査ログの定期確認
- バックアップの定期実行

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。

## コントリビューション

1. フォークする
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## サポート

質問や問題がある場合は、GitHub Issues でお知らせください。