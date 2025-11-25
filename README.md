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
