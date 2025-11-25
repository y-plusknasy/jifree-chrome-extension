# セキュリティと認証設計

## 認証方式 (簡易版)

スモールスタートのため、本格的なOAuth等は導入せず、以下の簡易的な認証とレート制限を組み合わせることで、不正利用とDoS攻撃を防ぐ。

### 1. クライアント認証 (共通鍵)
- **目的**: リクエストが正規のChrome拡張機能から送信されたものであることを簡易的に確認する。
- **実装**:
    - **Client**: 配布パッケージ内の設定ファイル（`.env` 等）に共通のハッシュ値（`SHARED_SECRET`）を保持。
    - **Server**: 環境変数に保持した正解のハッシュ値と照合。
    - **注意**: 難読化してもクライアントサイドにある以上、抽出は可能であるため、あくまで簡易的なフィルタリングとする。

### 2. ユーザー識別 (User ID)
- **目的**: ユーザーごとのレート制限を行うための一意な識別子。
- **実装**:
    - **Client**:
        - 初回起動時にランダムなUUID (v4) を生成。
        - `chrome.storage.local` に保存し、以降のリクエストで再利用する。
    - **Server**:
        - リクエストボディに含まれる `user_id` を検証に使用。

## APIリクエスト仕様

認証情報はPOSTリクエストのボディに含める。

```json
POST /analyze
{
  "html": "...",
  "selection": "...",
  "auth": {
    "shared_secret": "YOUR_SHARED_SECRET_HASH",
    "user_id": "generated-uuid-v4"
  }
}
```

## レート制限 (Rate Limiting)

ボットや悪意あるユーザーからのDoS攻撃を防ぐため、短期間の連続リクエストをブロックする。
コストとパフォーマンスを考慮し、外部データストア（Redis等）は使用せず、Cloud Functionsのインスタンスメモリを利用する。

- **データストア**: In-Memory (Go `sync.Map` or Cache Library)
- **ロジック**: Simple Expiration
- **制限ルール**:
    - **間隔**: 10秒
    - **挙動**:
        1. リクエスト受信時、メモリ内のキャッシュに `user_id` が存在するか確認。
        2. **存在する場合 (Cache Hit)**:
            - **ブロック**: HTTP 429 (Too Many Requests) を返す。
            - **TTL**: 更新しない（最初のブロックから10秒経過すれば解除）。
        3. **存在しない場合 (Cache Miss)**:
            - **許可**: 処理を続行。
            - **記録**: メモリに `user_id` を保存し、有効期限を10秒後に設定。

※ インスタンスメモリ依存のため、オートスケールで別インスタンスにリクエストが飛んだ場合はすり抜けるが、DoS攻撃への抑止力としては十分と判断する。

## CORS (Cross-Origin Resource Sharing)

ブラウザからの不正な呼び出しを防ぐため、CORSヘッダーを厳格に設定する。

- **Access-Control-Allow-Origin**: `chrome-extension://<EXTENSION_ID>`
- **Access-Control-Allow-Methods**: `POST, OPTIONS`
- **Access-Control-Allow-Headers**: `Content-Type`

※ `<EXTENSION_ID>` はビルド時または環境変数で設定。

### インフラ構成
- **Redis**: なし (In-Memory対応のため削除)
