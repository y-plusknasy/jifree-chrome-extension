# Jifree (ジフリー) - 漢字読み方支援 Chrome拡張機能

Webページ上でドラッグ選択した漢字にふりがな（ルビ）を表示する、子供向けの Chrome 拡張機能（モノレポ構成）。
技術スタック: 拡張機能 = TypeScript + Vite（Manifest V3）／バックエンド = Go 1.24 + kagome（形態素解析, IPADIC）on Google Cloud Functions。
開発はドキュメント駆動（`./docs/` を正とする）。

## 現在地
| マイルストーン | 状態 |
|---|---|
| 拡張機能・バックエンドの初期実装（v1.0.0） | 完了 |
| Chrome ウェブストア公開 | 完了 |
| 読み精度の改善（辞書 IPADIC → UniDic） | 未着手（方針決定済み） |

## 未完了タスク
- [ ] 誤読の実例 20〜50 文で評価セットを作り、IPADIC / UniDic の正解率を比較（[調査結果](./notes/research/2026-10-11-reading-accuracy-improvement.md)）
- [ ] 比較結果を踏まえて UniDic へ切替し、方針を ADR 化
- [ ] ドキュメント配置の再編: GitHub Pages 用ページを専用ディレクトリ（`public` 等）へ移し、`.md` 系ドキュメントを最終的に `./docs/` へ集約（`./notes/` もここへ移す）

## 次回への引き継ぎ
- 追加したい機能の要望は、ユーザーが別セッションで提示予定
- UniDic は分割単位が細かく読みのフィールド位置も異なるため、`findTargetTokens` の照合とフロントのルビ付与への影響を要確認
- `./docs/` は GitHub Pages（main:/docs、公開リポジトリ）の公開元。プライバシーポリシー等の URL はストア掲載から参照されている可能性が高く、移動時はリンク切れに注意
- 本番の秘密値は `./.env` と `./extension/.env.production`（Git 管理外）。デプロイは `make deploy`、ストア提出用 ZIP は `make package-extension`
- 自動テストは未整備

## 主要ファイル
- [README](./README.md)（ローカル動作確認手順）
- [Makefile](./Makefile)（デプロイ・ビルド・パッケージング）
- [拡張機能](./extension/)・[manifest.json](./extension/manifest.json)
- [バックエンド](./backend/)（エントリーポイント: [function.go](./backend/function.go)、読み算出: [analyzer.go](./backend/internal/morph/analyzer.go)）
- [設計ドキュメント](./docs/architecture/)・[API 仕様](./docs/api/openapi.yaml)
- [公開ページ](./docs/index.md)・[プライバシーポリシー](./docs/privacy-policy.md)
- [調査結果](./notes/research/)・[セッションメモ](./notes/memo/)（Pages 公開対象外の暫定置き場）
- ADR（設計判断の記録）: 未作成

## ADR索引
<!-- - [ADR 0001: タイトル](./docs/adr/adr-0001-xxx.md) — 一言要約 -->
