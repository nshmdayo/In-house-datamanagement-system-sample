# プライベートブロックチェーン社内データ管理システム - 開発ガイドライン

## プロジェクト概要

このプロジェクトは、プライベートブロックチェーン技術を活用した企業向け社内データ管理システムです。データの改ざん防止、透明性の確保、アクセス制御を実現し、企業の重要なデータを安全に管理します。

## 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **ブロックチェーン**: Hyperledger Fabric または独自実装
- **データベース**: PostgreSQL
- **API**: REST API (Gin Framework)
- **認証**: JWT + RBAC (Role-Based Access Control)
- **暗号化**: AES-256, RSA

### フロントエンド（将来実装）
- **フレームワーク**: React/Next.js
- **状態管理**: Redux Toolkit
- **UI**: Material-UI

## システム要件

### 機能要件

#### 1. ユーザー管理
- ユーザー登録・認証（管理者による承認制）
- 役割ベースのアクセス制御（Admin, Manager, Employee, Guest）
- プロファイル管理
- ログイン履歴の追跡

#### 2. データ管理
- ドキュメントの登録・更新・削除
- バージョン管理
- メタデータ管理（作成者、更新日時、カテゴリ等）
- ファイルの暗号化保存

#### 3. ブロックチェーン機能
- データ操作の全履歴をブロックチェーンに記録
- 改ざん検知機能
- データの整合性検証
- トランザクションの追跡

#### 4. アクセス制御
- ファイル/フォルダ単位でのアクセス権限設定
- 部門別のデータアクセス制御
- 機密レベルに応じたアクセス制限

#### 5. 監査機能
- 全ての操作ログの記録
- アクセスログの可視化
- コンプライアンスレポート生成

#### 6. 検索・フィルタリング
- 全文検索
- メタデータ検索
- 高度なフィルタリング機能

### 非機能要件

#### 1. セキュリティ
- データの暗号化（保存時・転送時）
- 多要素認証（MFA）
- セキュリティログの監視
- 定期的なセキュリティ監査

#### 2. パフォーマンス
- API応答時間: 95%以上のリクエストが500ms以内
- 同時接続数: 1000ユーザー
- データベース応答時間: 100ms以内

#### 3. 可用性
- システム稼働率: 99.9%以上
- 自動バックアップ（日次・週次）
- 災害復旧計画

#### 4. スケーラビリティ
- 水平スケーリング対応
- マイクロサービスアーキテクチャ
- コンテナ化（Docker）

## アーキテクチャ設計

### システム構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Backend       │
│   (React)       │◄──►│   (Nginx)       │◄──►│   (Go)          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                         │
                         ┌─────────────────┐            │
                         │  Blockchain     │◄───────────┤
                         │  Network        │            │
                         └─────────────────┘            │
                                                         │
                         ┌─────────────────┐            │
                         │  PostgreSQL     │◄───────────┘
                         │  Database       │
                         └─────────────────┘
```

### ディレクトリ構造

```
├── cmd/                    # アプリケーションエントリーポイント
│   └── server/
├── internal/               # 内部パッケージ
│   ├── api/               # API関連
│   │   ├── handlers/      # HTTPハンドラー
│   │   ├── middleware/    # ミドルウェア
│   │   └── routes/        # ルーティング
│   ├── blockchain/        # ブロックチェーン関連
│   │   ├── block/         # ブロック構造
│   │   ├── chain/         # チェーン管理
│   │   ├── consensus/     # コンセンサス
│   │   └── transaction/   # トランザクション
│   ├── config/            # 設定管理
│   ├── database/          # データベース関連
│   │   ├── migrations/    # マイグレーション
│   │   └── models/        # データモデル
│   ├── security/          # セキュリティ関連
│   │   ├── auth/          # 認証
│   │   ├── crypto/        # 暗号化
│   │   └── rbac/          # 役割ベースアクセス制御
│   └── services/          # ビジネスロジック
├── pkg/                   # 外部パッケージ
├── scripts/               # スクリプト
├── deployments/           # デプロイメント設定
├── docs/                  # ドキュメント
└── tests/                 # テスト
```

## データモデル

### 主要エンティティ

#### User（ユーザー）
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Role      Role      `json:"role"`
    IsActive  bool      `json:"is_active" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### Document（ドキュメント）
```go
type Document struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Content     string    `json:"content"`
    FileHash    string    `json:"file_hash" gorm:"unique"`
    Category    string    `json:"category"`
    AccessLevel int       `json:"access_level"`
    CreatedBy   uint      `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### BlockchainRecord（ブロックチェーン記録）
```go
type BlockchainRecord struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    TransactionID string    `json:"transaction_id" gorm:"unique"`
    BlockHash     string    `json:"block_hash"`
    DocumentID    uint      `json:"document_id"`
    Action        string    `json:"action"`
    UserID        uint      `json:"user_id"`
    Timestamp     time.Time `json:"timestamp"`
}
```

## API設計

### 認証エンドポイント
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークン更新
- `POST /api/v1/auth/logout` - ログアウト

### ユーザー管理
- `GET /api/v1/users` - ユーザー一覧取得
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users/{id}` - ユーザー詳細取得
- `PUT /api/v1/users/{id}` - ユーザー更新
- `DELETE /api/v1/users/{id}` - ユーザー削除

### ドキュメント管理
- `GET /api/v1/documents` - ドキュメント一覧取得
- `POST /api/v1/documents` - ドキュメント作成
- `GET /api/v1/documents/{id}` - ドキュメント詳細取得
- `PUT /api/v1/documents/{id}` - ドキュメント更新
- `DELETE /api/v1/documents/{id}` - ドキュメント削除

### ブロックチェーン
- `GET /api/v1/blockchain/blocks` - ブロック一覧取得
- `GET /api/v1/blockchain/transactions/{id}` - トランザクション詳細
- `POST /api/v1/blockchain/verify` - データ整合性検証

## セキュリティガイドライン

### 暗号化
- データベースの機密データはAES-256で暗号化
- パスワードはbcryptでハッシュ化
- API通信はTLS 1.3を使用

### 認証・認可
- JWT トークンの有効期限は15分
- リフレッシュトークンの有効期限は7日
- RBAC による細かなアクセス制御

### 入力検証
- 全ての入力データのバリデーション
- SQLインジェクション対策
- XSS対策

## テスト戦略

### 単体テスト
- カバレッジ 80% 以上
- モック使用によるテスト分離

### 統合テスト
- API エンドポイントのテスト
- データベース連携テスト

### セキュリティテスト
- 脆弱性スキャン
- ペネトレーションテスト

## 開発規約

### コードスタイル
- gofmt によるコード整形
- golint によるリント
- ネーミング規約はGo標準に従う

### コミット規約
- Conventional Commits 形式を採用
- feat: 新機能
- fix: バグ修正
- docs: ドキュメント
- test: テスト

### パッケージ管理
- Go Modules を使用
- 依存関係の定期的な更新

## デプロイメント

### 環境
- Development
- Staging  
- Production

### CI/CD
- GitHub Actions を使用
- 自動テスト実行
- Docker イメージビルド
- K8s への自動デプロイ

---

このガイドラインに従って、安全で信頼性の高い社内データ管理システムを構築してください。
