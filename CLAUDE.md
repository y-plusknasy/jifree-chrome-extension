# Jifree (ジフリー) - 漢字読み方支援 Chrome拡張機能

Webページ上でドラッグ選択した漢字にふりがな（ルビ）を表示する、子供向けの Chrome 拡張機能（モノレポ構成）。
技術スタック: 拡張機能 = TypeScript + Vite（Manifest V3）／バックエンド = Go 1.24 + kagome（形態素解析）on Google Cloud Functions。
開発はドキュメント駆動（`./docs/` を正とする）。

## 現在地
| マイルストーン | 状態 |
|---|---|
| 拡張機能・バックエンドの初期実装 | 完了 |
| Chrome ウェブストア公開 | 進行中（レビュー指摘による manifest 調整済み） |

## 未完了タスク
- [ ] Chrome ウェブストア公開の完了

## 次回への引き継ぎ
- `.gitignore` の未コミット変更あり（`screenshot.png` の除外、`/gi` による秘匿情報・ログ除外の追記）
- 本番の秘密値は `./.env` と `./extension/.env.production`（どちらも Git 管理外）。デプロイは `make deploy`、ストア提出用 ZIP は `make package-extension`

## 主要ファイル
- [README](./README.md)（ローカル動作確認手順）
- [Makefile](./Makefile)（デプロイ・ビルド・パッケージング）
- [拡張機能](./extension/)・[manifest.json](./extension/manifest.json)
- [バックエンド](./backend/)（エントリーポイント: [function.go](./backend/function.go)）
- [設計ドキュメント](./docs/architecture/)・[API 仕様](./docs/api/openapi.yaml)
- [公開ページ](./docs/index.md)・[プライバシーポリシー](./docs/privacy-policy.md)
- [ADR](./docs/adr/)（設計判断の記録。未作成）

## ADR索引
<!-- - [ADR 0001: タイトル](./docs/adr/adr-0001-xxx.md) — 一言要約 -->
