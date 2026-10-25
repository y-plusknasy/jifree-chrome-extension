# Jifree (ジフリー) - 漢字読み方支援 Chrome拡張機能

Webページ上でドラッグ選択した漢字にふりがな（ルビ）を表示する、子供向けの Chrome 拡張機能（モノレポ構成）。
技術スタック: 拡張機能 = TypeScript + Vite（Manifest V3）／バックエンド = Go 1.24 + kagome（形態素解析, IPADIC）on Google Cloud Functions。
開発はドキュメント駆動。
ただし既存のドキュメント（初期に GitHub Copilot で作成）とコードの照合は実施済みで、食い違いは [食い違い一覧](./notes/doc-code-discrepancies.md) に記録した（どちらが正しいかは未判断）。
食い違いはアップデートで触る範囲だけをその都度直す方針のため、触っていない箇所は古いまま残る。ドキュメントを根拠に使う前に一覧を確認し、コードで裏付けを取ること。

## 現在地
| マイルストーン | 状態 |
|---|---|
| 拡張機能・バックエンドの初期実装（v1.0.0） | 完了 |
| Chrome ウェブストア公開 | 完了 |
| ドキュメントとコードの照合（食い違いの記録） | 完了 |
| 次回アップデート（辞書の見直し・公開ページと実装の統一） | 検討中（辞書の比較・追加機能の要望の提示と調査は完了、方向と採否は未決定） |

## 未完了タスク
### 次回アップデート（以下の順に進める）
- 目的: ①辞書まわりを見直して読みの精度を上げる（UniDic への単純な切り替えは実測で効果なし＝[比較結果](./notes/research/2026-10-24-ipadic-vs-unidic.md)） ②公開ページと実装の食い違いを直して統一する
- 進め方: 食い違いの修正は独立したタスクにしない。アップデートで触る範囲だけを、矛盾がなくなるようその都度直す（範囲外は一覧で「未対応」のまま残す）
- [ ] 1. 今回のアップデートで加える変更点を決める（追加機能の要望を含む。要望は Git 管理外のローカルメモ `./notes/memo/2026-10-25-feature-ideas.md` に記録済み・採否は未決定）
  - [ ] 辞書の方向を決める（A: IPADIC＋ユーザー辞書 / B: 商用 LLM API 併用 / C: UniDic 採用。[比較結果](./notes/research/2026-10-24-ipadic-vs-unidic.md) の 6 章）
- [ ] 2. 変更で触るファイルの範囲を決め、範囲内の食い違いを[一覧](./notes/doc-code-discrepancies.md)から抜き出す
- [ ] 3. 2 の食い違いの解消と追加実装のプランを作り、実装する（辞書の方針は ADR 化。根拠は[比較結果](./notes/research/2026-10-24-ipadic-vs-unidic.md)と[前回の調査](./notes/research/2026-10-11-reading-accuracy-improvement.md)）
- [ ] 4. 範囲内のドキュメントを実装に合わせて更新し、一覧の該当項目の「状態」を更新する。ドキュメント配置の再編（GitHub Pages 用ページを専用ディレクトリ〔`public` 等〕へ移し、`.md` 系ドキュメントを `./docs/` へ集約。`./notes/` もここへ移す）もここで行う
- [ ] 5. 最新の内容を公開ページとストア掲載に反映する（データ開示とプライバシーポリシーの整合、「データ収集なし」の言い回し、相互リンク、配置変更によるリンク切れの確認）。反映までは現行版の掲載を据え置く

### その他
- [ ] v1.0.0 の Git タグと GitHub リリースを作成（提出版のコミット特定から。以後はストア提出ごとにタグ付け）（[公開先・リリース情報](./notes/release-info.md)）

## 次回への引き継ぎ
- 追加したい機能の要望は 2026-10-25 に提示済み。機能を大きく変えるか迷っている段階で、採否は未決定。内容は非公開とするため Git 管理外のローカルメモ `./notes/memo/2026-10-25-feature-ideas.md` にのみ記録（このPCにしかない）。同日の調査結果も Git 管理外の `./notes/private/research/2026-10-25-feature-feasibility.md` にある（`./notes/private/` は非公開資料の置き場）。次に決めることはメモの末尾を参照
- 辞書比較（例文 44 問）は IPADIC 35・UniDic（読み補正あり）35 の同点。kagome 標準の UniDic は仮名形の列がなくメモリ約 4 倍。例文は自作中心なので、方向を決める前に実際の誤読例を加えて再測定すると判断材料が増える（比較プログラムはリポジトリに保存していない。手法は比較結果の 4 章）
- `./docs/` は GitHub Pages（main:/docs、公開リポジトリ）の公開元。プライバシーポリシーの URL はストア掲載から参照されている（確認済み）ため、移動時はリンク切れに注意
- リポジトリは公開設定のため、`./notes/` も GitHub 上では誰でも閲覧できる（Pages に載らないだけ）。個人情報は書かない。例外として `./notes/private/` と上記の要望メモは `.gitignore` で除外しており公開されない（非公開にしたい資料はここへ）
- 本番の秘密値は `./.env` と `./extension/.env.production`（Git 管理外）。デプロイは `make deploy`、ストア提出用 ZIP は `make package-extension`
- 自動テストは未整備

## 主要ファイル
- [README](./README.md)（ローカル動作確認手順）
- [Makefile](./Makefile)（デプロイ・ビルド・パッケージング）
- [拡張機能](./extension/)・[manifest.json](./extension/manifest.json)
- [バックエンド](./backend/)（エントリーポイント: [function.go](./backend/function.go)、読み算出: [analyzer.go](./backend/internal/morph/analyzer.go)）
- [設計ドキュメント](./docs/architecture/)・[API 仕様](./docs/api/openapi.yaml)（照合済み。食い違いは一覧を参照）
- [ドキュメントとコードの食い違い一覧](./notes/doc-code-discrepancies.md)
- [公開ページ](./docs/index.md)・[プライバシーポリシー](./docs/privacy-policy.md)
- [公開先・リリース情報](./notes/release-info.md)（ストア掲載 URL・GitHub Pages・リリース状況）
- [調査結果](./notes/research/)・[セッションメモ](./notes/memo/)（Pages 公開対象外の暫定置き場）
- ADR（設計判断の記録）: 未作成

## ADR索引
<!-- - [ADR 0001: タイトル](./docs/adr/adr-0001-xxx.md) — 一言要約 -->
