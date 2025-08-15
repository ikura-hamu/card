# card プロジェクト - AI エージェント向けガイド

## プロジェクト概要

これは Go + Bubble Tea を使用したターミナルベースの個人プロフィール名刺アプリケーションです。GitHub のプロフィール情報、リポジトリ一覧、アイコンを表示する TUI（Terminal User Interface）として動作します。

## アーキテクチャの理解

### コア設計パターン

- **タブベースアーキテクチャ**: `internal/tabs/tabs.go` がメインの UI コントローラー
- **Model-View-Update (MVU) パターン**: Bubble Tea の Elm アーキテクチャに従う
- **タブインターフェース**: すべてのタブは `Tab` インターフェースを実装

```go
type Tab interface {
    tea.Model
    Name() string
}
```

### モジュール構成

```
internal/
├── tabs/          # タブマネージャー（メインコントローラー）
├── about/         # GitHubプロフィールREADME表示
├── repo/          # リポジトリ一覧とブラウザ連携
├── icon/          # GitHub アバター画像のASCII変換
└── common/        # 共通ユーティリティ
    ├── merrors/   # エラーハンドリング用コマンド生成
    └── size/      # 画面サイズ管理
```

## 重要な開発パターン

### エラーハンドリング

- エラーは `merrors.NewCmd(err)` を使用して Bubble Tea コマンドとして返す
- エラーメッセージは `TabsManager.Update()` で捕捉され、`tea.Quit` でアプリを終了

### サイズ管理

- `tea.WindowSizeMsg` はタブマネージャーでコンテンツエリアサイズに変換
- 各タブは変換されたサイズを受け取り、ボーダーやパディングを考慮した描画を行う
- `size.Size` 構造体で幅と高さを管理

### 非同期データ取得

- HTTP リクエストは `tea.Cmd` として実行（例: `fetchReadme`, `fetchRepositories`）
- GitHub API レート制限を適切に処理（`reposRateLimitMsg`）

### スタイリング

- Lipgloss を使用した一貫したスタイリング
- タブヘッダー: アクティブ/非アクティブ状態の視覚的区別
- コンテンツエリア: 統一されたボーダーと背景色

## 開発ワークフロー

### ビルドとデバッグ

- タスク: "g++ デバッグ用コンパイル" が利用可能
- 実行: `go run main.go` でアプリケーション起動
- 依存関係: `go mod tidy` で管理

### 新しいタブの追加

1. `internal/` 下に新しいパッケージを作成
2. `Tab` インターフェースを実装
3. `main.go` の `tabs.NewTabsManager()` に追加

### テンポラリファイル管理

- アイコンは一時ファイルとして保存
- `icon.CleanupTempIcons()` で終了時にクリーンアップ
- PID ベースのプレフィックスで他プロセスとの競合を回避

## 外部依存関係

- **Bubble Tea**: TUI フレームワーク
- **Lipgloss**: ターミナルスタイリング
- **Glamour**: マークダウンレンダリング
- **go-github**: GitHub API クライアント
- **ascii-image-converter**: 画像の ASCII 変換
- **browser**: URL をブラウザで開く

## 設定とカスタマイズ

- GitHub ユーザー名: `internal/repo/github.go` の `githubUserName` 定数
- アイコン URL: `internal/icon/icon.go` の `iconURL` 定数
- スタイル設定: 各モジュールの末尾にある lipgloss スタイル変数
