# Jifree (Kanji Reading Helper)

子供向けの、漢字の読み方を教えるChrome拡張機能プロジェクトです。
モノレポ構成で管理されています。

## プロジェクト構成

- **extension/**: Chrome拡張機能のソースコード。ユーザーのドラッグ操作をトリガーにHTML要素を取得します。
- **backend/**: Google Cloud Functions (Go runtime)。形態素分析を行い、読み仮名を返します。
- **docs/**: 設計書、仕様書、アーキテクチャ図。

## 開発の進め方

ドキュメント駆動開発 (DDD) を採用しています。
`docs/` 以下のドキュメントを正として開発を進めます。

## ローカルでの動作確認

バックエンドサーバーを起動した状態で、以下のコマンドでAPIの動作確認ができます。
※ `docker-compose.yml` に `SHARED_SECRET=test-secret` が設定されている前提です。

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -H "Origin: chrome-extension://dummy-id" \
  -d '{
    "html": "<p>私の妻は最高にエキゾチックで<strong>天女</strong>のような女性です。</p>",
    "selection": "天女",
    "prefix": "エキゾチックで",
    "suffix": "のような",
    "auth": {
      "shared_secret": "test-secret",
      "user_id": "test-user-001"
    }
  }'
```
