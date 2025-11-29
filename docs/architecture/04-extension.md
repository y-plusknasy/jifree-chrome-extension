# Chrome拡張機能 設計 (Client)

## 技術スタック (V1)

開発効率とビルドの最適化のため、Viteを採用する。
フレームワーク（React等）は使用せず、Vanilla JSで実装する。

- **言語**: TypeScript
- **ビルド/バンドル**: Vite
- **UI**: Vanilla JS (DOM API)
- **Manifest**: Version 3

## ディレクトリ構成

```text
extension/
├── manifest.json
├── package.json       # 依存関係定義
├── tsconfig.json
├── vite.config.ts     # ビルド設定 (Multi-entry)
├── src/
│   ├── background.ts      # Service Worker (API通信、認証付与)
│   ├── content.ts         # DOM操作 (イベント検知、UI表示)
│   ├── types.ts           # 型定義 (APIレスポンス等)
│   └── utils/             # 共通ロジック (テキスト抽出等)
└── dist/                  # ビルド出力先 (Chrome読み込み用)
    ├── background.js
    ├── content.js
    └── assets/            # CSS等 (必要に応じて)
```

## コンポーネント設計

### 1. Content Script (`content.ts`)
Webページ上で直接動作するスクリプト。

- **イベント監視**:
    - `mouseup`: テキスト選択の完了を検知。
    - `mousedown`: ポップアップ外クリックによる閉じる動作を検知。
- **コンテキスト抽出**:
    - `window.getSelection()` から `Range` オブジェクトを取得。
    - 選択範囲のテキスト (`selection`) に加え、前後の文脈 (`prefix`, `suffix`) と親要素のHTML (`html`) を抽出する。
    - 文脈マッチング精度向上のため、親要素が短い場合はさらに親へ遡るロジックを含む。
- **通信**:
    - `chrome.runtime.sendMessage` を使用して Background Script に解析リクエストを送信。
- **UI表示**:
    - APIレスポンス（読み仮名トークン）を受け取り、選択範囲の下部にポップアップを表示する。
    - スタイルは `style` タグを動的に注入し、クラス名プレフィックス (`jifree-`) を使用して既存サイトとの衝突を避ける。

### 2. Background Script (`background.ts`)
拡張機能のバックグラウンドで動作する Service Worker。

- **メッセージ受信**:
    - Content Script からの解析リクエストを受信。
- **認証処理**:
    - `chrome.storage.local` からユーザーID (UUID) を取得（なければ生成・保存）。
    - 設定ファイル（`src/utils/config.ts`）から共通鍵 (`SHARED_SECRET`) を取得。
- **API通信**:
    - バックエンドAPI (`POST /analyze`) を実行。
    - 認証情報 (`auth`) をリクエストボディに付与。
- **レスポンス返却**:
    - APIの結果を Content Script に返す。

## データフロー

1. **User**: テキストをドラッグして選択。
2. **Content**: `mouseup` 検知 -> `Selection`, `Prefix`, `Suffix`, `HTML` 抽出 -> `sendMessage`。
3. **Background**: `onMessage` 受信 -> `UserID`, `Secret` 付与 -> `fetch(API)`。
4. **API (Backend)**: 解析実行 -> JSON返却。
5. **Background**: JSON受信 -> `sendResponse`。
6. **Content**: JSON受信 -> ポップアップ生成 -> DOMに追加。
