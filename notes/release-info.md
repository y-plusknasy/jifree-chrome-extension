# 公開先・リリース情報

Jifree の外部公開先（GitHub / GitHub Pages / Chrome ウェブストア）の現状記録。
**最終確認日: 2026-10-24**（各値は確認日時点。変更時はこのファイルを更新する）

> 注意: リポジトリは公開設定のため、`./notes/` も GitHub 上では誰でも閲覧できる（Pages に載らないだけ）。連絡先などの個人情報はここに書かない。

## GitHub リポジトリ
| 項目 | 内容 |
|---|---|
| URL | https://github.com/y-plusknasy/jifree-chrome-extension |
| 公開範囲 | Public |
| 作成日 | 2025-11-25 |
| About 欄 | 説明 `Jifree - Kanji Reading Helper`、ホームページ欄に GitHub Pages の URL を登録済み |

## GitHub Pages
| 項目 | 内容 |
|---|---|
| 公開 URL | https://y-plusknasy.github.io/jifree-chrome-extension/ |
| 公開元 | `main` ブランチの `/docs`（従来型ビルド、HTTPS 強制、独自ドメインなし） |
| 最終ビルド | 2026-10-11（コミット `85b47e8`） |

| ページ | 公開 URL | 元ファイル |
|---|---|---|
| トップ（紹介・使い方デモ） | https://y-plusknasy.github.io/jifree-chrome-extension/ | [./docs/index.md](../docs/index.md) |
| プライバシーポリシー（最終更新日 2025-12-01） | https://y-plusknasy.github.io/jifree-chrome-extension/privacy-policy | [./docs/privacy-policy.md](../docs/privacy-policy.md) |

- 公開ページの構成・リンクは元ファイルと一致していることを確認済み
- `./docs/` 配下の設計文書（`architecture/`・`api/`）も同じ仕組みで公開対象になる
- **プライバシーポリシーの URL は Chrome ウェブストア掲載から参照されている（確認済み）**。移動・改名するとストア掲載がリンク切れになる

## Chrome ウェブストア
| 項目 | 内容 |
|---|---|
| 掲載 URL | https://chromewebstore.google.com/detail/jifree-kanji-reading-help/nlbfdkbnbnpfoiiobgjihjbodejoomch |
| 拡張機能 ID | `nlbfdkbnbnpfoiiobgjihjbodejoomch` |
| 名前 | Jifree - Kanji Reading Helper |
| 概要 | 「Webページ上の選択した漢字に読み仮名を表示します。」（`manifest.json` の `description` と同一） |
| 公開中バージョン | 1.0.0 |
| 最終更新日 | 2025-12-07 |
| サイズ / 言語 | 18.69KiB / 日本語 |
| カテゴリ | 教育 |
| ユーザー数 / 評価 | 79 ユーザー / レビューなし |
| 提供元表示 | +knasy（連絡先メールは掲載ページ参照。ウェブサイト欄は未設定） |
| プライバシーポリシー | GitHub Pages の `/privacy-policy` |
| データ取り扱いの開示 | 「ユーザーのアクティビティ」「ウェブサイトのコンテンツ」の2種類 |
| 公開範囲 | 一般公開と判断（未ログインで閲覧でき、ユーザー数も表示されるため。管理画面では未確認） |

- 管理画面: Chrome ウェブストア デベロッパー ダッシュボード（https://chrome.google.com/webstore/devconsole ）
- 提出用 ZIP は `make package-extension` で作成（[Makefile](../Makefile)）

## GitHub リリース
- **リリースは 0 件、Git タグも 0 件**（2026-10-24 時点、API で確認）
- v1.0.0 の痕跡は、コミット `be195d7`（2025-11-30 `feat: Release v1.0.0 - Initial deployment ready`）と、`manifest.json`・`package.json` の `version: 1.0.0` のみ
- ストアに提出したコミットは記録がない。ストア最終更新（2025-12-07）直前のコード変更は `e5b2731`（2025-12-03 `Fix: Adjust manifest.json per review feedback`）であり、これが提出版と**推定**される（未検証）

## 気づいた点（現行版は据え置き、次のアップデート版リリース時にまとめて見直す）
1. ストアのデータ開示は2種類だが、プライバシーポリシーは匿名識別子（UUID）も記載している。開示の要否を確認する
2. 公開ページの「余計な広告やデータ収集は一切ありません」という言い回しが、ストアの開示内容と食い違って見える
3. 公開ページ → ストア、ストア → 公開ページ（ウェブサイト欄）の相互リンクがない
4. 提出版とコミットの対応が残っていない（今後は提出ごとにタグを付ける）
