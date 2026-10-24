# ドキュメントとコードの食い違い一覧

照合日: 2026-10-24 ／ 照合時点のコミット: `0e0693c`

## この一覧の扱い
- ここに挙げた項目は、**ドキュメントとコードのどちらが正しいかを判断していない**。記録のみで、どちらも修正していない。
- 該当箇所に手を加えることになったときに、議論して「どちらに合わせるか」を決め、対応する。
- アップデートで触らない箇所は「未対応」のまま残す（未対応が残ること自体は問題としない）。
- 対応したら「状態」を更新する（例: `対応済み（コードに合わせてドキュメント修正, コミット xxx）`）。不要になった項目は理由を添えて「対応不要」にする。

## 照合範囲
| ドキュメント | 照合したコード・設定 |
|---|---|
| `./docs/api/openapi.yaml` | `./backend/function.go`・`./backend/internal/api/types.go`・`./extension/src/types.ts`・`./extension/src/background.ts` |
| `./docs/architecture/01〜04` | `./backend/` 全体・`./extension/` 全体・`./extension/manifest.json`・`./Makefile` |
| `./README.md` | `./Makefile`・`./docker-compose.yml`・`./.env.example`・`./extension/.env.development` |
| `./docs/index.md`・`./docs/privacy-policy.md`（公開ページ） | 拡張機能が送信・保存するデータ、バックエンドの処理 |

ドキュメント同士の食い違いも含めている。Git 管理外の本番設定（`./.env`・`./extension/.env.production`）は開いていないため、本番値に依存する点は「要確認事項」に分けた。

重要度の目安: **高** = 公開ページ・ストアの開示内容に関わる ／ **中** = 仕様や動作を読み違える恐れがある ／ **低** = 表記ゆれ・省略・軽微な記述差

## 一覧
| ID | 対象 | 概要 | 重要度 | 状態 |
|---|---|---|---|---|
| D-01 | プライバシーポリシー | 送信データは「選択テキスト」と記載。実際は周辺要素の HTML なども送信している | 高 | 未対応 |
| D-02 | 公開ページ | 「データ収集は一切なし」と記載。プライバシーポリシーとコードではデータを送信している | 高 | 未対応 |
| D-03 | 設計書・プライバシーポリシー | ユーザー ID の生成を「初回起動時」と記載。実際は初回の解析リクエスト時 | 中 | 未対応 |
| D-04 | 公開ページ・設計書 | 「漢字の上にふりがな」と記載。実際は全単語にカタカナのルビを付けたポップアップ | 中 | 未対応 |
| D-05 | API 仕様・設計書 | API のパスを `/analyze` と記載。コードにはパスの振り分けがない | 中 | 未対応 |
| D-06 | データフロー図 | レスポンスのキーが `readings`。コードと API 仕様は `tokens` | 中 | 未対応 |
| D-07 | セキュリティ設計・API 仕様 | レート制限（連続リクエストの制限）の秒数と解除の起点が、記述どうしでもコードとも食い違う | 中 | 未対応 |
| D-08 | API 仕様 | 必須パラメータが欠けると 400 と記載。コードは必須チェックをしていない | 中 | 未対応 |
| D-09 | セキュリティ設計 | 呼び出し元（Origin）を厳格に制限と記載。コードは設定がなければ制限しない | 中 | 未対応 |
| D-10 | README | ローカル確認の前提が docker compose に限られる。`make run-backend` では手順どおりに動かない | 中 | 未対応 |
| D-11 | API 仕様・README | `user_id` は UUID 形式と記載。コードは形式を検証せず、README の例も UUID ではない | 低 | 未対応 |
| D-12 | セキュリティ設計 | 共通鍵を「ハッシュ値」と記載。コードは文字列をそのまま比較 | 低 | 未対応 |
| D-13 | 拡張機能設計 | ディレクトリ構成の説明が実際のファイルと一部違う | 低 | 未対応 |
| D-14 | データフロー図・セキュリティ設計 | リクエスト例に `prefix`／`suffix` がない | 低 | 未対応 |
| D-15 | 概要設計 | 操作を「ドラッグ＆ドロップ」と記載。実際はドラッグによるテキスト選択 | 低 | 未対応 |
| D-16 | README | 「docs/ を正とする」と記載。CLAUDE.md の現状認識（整合は未保証）と食い違う | 低 | 未対応 |
| D-17 | API 仕様 | 読みが空文字になる場合があると記載。コードは空文字を返さない | 低 | 未対応 |
| D-18 | API 仕様 | 文脈が一致しないときの代替処理（選択テキストだけを解析）が書かれていない | 低 | 未対応 |
| C-01〜03 | コード内コメント | コメントと処理内容が食い違う箇所 | 低 | 未対応 |

## 詳細

### D-01 送信データの範囲（高）
- **ドキュメント**: `./docs/privacy-policy.md` L11-14 では、送信するのは「ユーザーが Web ページ上でドラッグして選択したテキストデータ」とだけ書いている。L21-24 では、ユーザー ID はローカルに保存するとだけ書いている。
- **コード**:
  - 選択テキストに加えて、次のものも送信している（`./extension/src/content.ts` L20, L60-108、`./extension/src/background.ts` L40-49）。
    - 直前・直後の各 10 文字（`prefix`／`suffix`）
    - 選択範囲を含む要素の HTML（`outerHTML`。タグや属性を含む。本文が 50 文字未満なら、`section`／`body` に達するまで親要素へさかのぼる）
  - ユーザー ID も毎回サーバーへ送信している（`background.ts` L45-48）。サーバー側では、レート制限用にインスタンス（サーバーの実行単位）が動いている間メモリに残る。期限切れの記録も削除しない（`./backend/internal/morph/ratelimit/limiter.go` L28-46）。
  - すべてのページで動作し（`./extension/manifest.json` L21 `<all_urls>`）、空でない選択なら漢字を含まなくても送信する（`content.ts` L5-13）。
- **論点**: ストアのデータ開示とも関係する。[CLAUDE.md](../CLAUDE.md) の「次回アップデート」の手順5（公開ページとストア掲載への反映）と合わせて扱う。

### D-02 「データ収集なし」の表現（高）
- **ドキュメント**: `./docs/index.md` L17「余計な広告やデータ収集は一切ありません」。
- **食い違う先**: `./docs/privacy-policy.md` 第 1 章は「収集・使用するデータ」を列挙している。コードも D-01 のとおりデータを送信している。
- **論点**: CLAUDE.md の「次回アップデート」の手順5（「データ収集なし」の言い回しの見直し）で扱う。

### D-03 ユーザー ID の生成タイミング（中）
- **ドキュメント**: `./docs/architecture/03-security.md` L18、`./docs/privacy-policy.md` L24 は「初回起動時に生成」。
- **コード**: 解析リクエストのたびに `getUserId()` を呼び、保存済みの ID がなければその時点で生成する（`./extension/src/background.ts` L10-20, L38）。インストール時や起動時の処理はない。
- **補足**: `./docs/architecture/04-extension.md` L56（「なければ生成・保存」）はコードと一致する。ドキュメント同士でも食い違っている。

### D-04 ふりがなの表示方法（中）
- **ドキュメント**:
  - `./docs/index.md` L11「自動的に漢字の上にふりがなが表示されます」。L34 は読みをひらがなで例示している（「はは」「ふぼ」）。
  - `./docs/architecture/01-overview.md` L20「漢字に対する読み仮名を特定」。
- **コード**:
  - 選択範囲の下にポップアップを出し、その中で全トークン（単語）にルビを付ける。ひらがなや記号のトークンも対象になる（`./extension/src/content.ts` L168-187）。
  - 読みは辞書の値（カタカナ）をそのまま使い、ひらがなには変換していない（`./backend/internal/morph/analyzer.go` L42-45、`content.ts` L181-182）。
- **補足**: ポップアップで表示する点は `./docs/architecture/04-extension.md` L47 と一致する。

### D-05 API のパス（中）
- **ドキュメント**: `./docs/api/openapi.yaml` L13、`./docs/architecture/02-data-flow.mermaid` L10、`./docs/architecture/03-security.md` L28、`./docs/architecture/04-extension.md` L59 はいずれも `POST /analyze`。
- **コード**:
  - バックエンドは関数を 1 つ登録しているだけで、パスによる振り分けがない（`./backend/function.go` L46）。
  - 拡張機能は設定値 `VITE_API_ENDPOINT` の URL をそのまま呼び出している（`./extension/src/background.ts` L51、`./extension/.env.development` L1 は `http://localhost:8080`）。
  - `./README.md` L23 もルート（`/`）へ送信している。
- **関連**: 本番の URL は要確認事項 Q-01。

### D-06 レスポンスのキー名（中）
- **ドキュメント**: `./docs/architecture/02-data-flow.mermaid` L21 は `JSON { readings: [...] }`。
- **コード**: キー名は `tokens`（`./backend/internal/api/types.go` L20、`./extension/src/types.ts` L13）。`./docs/api/openapi.yaml` L75 も `tokens`。

### D-07 レート制限の秒数と解除の起点（中）
- **ドキュメント**:
  - `./docs/architecture/03-security.md` は、L47 で間隔を「5 秒」としながら、L52 と L55 では「10 秒」としている（ドキュメント内の食い違い）。解除の起点は「最初のブロックから」（L52）。
  - `./docs/api/openapi.yaml` L100 は「5 秒以内に再リクエスト」と固定値で書いている。
- **コード**:
  - 既定値は 5 秒で、環境変数 `RATE_LIMIT_WINDOW_SEC` で変更できる（`./backend/function.go` L30-37。本番は `./Makefile` L14 で設定）。
  - 期限は**最後に許可されたリクエストから**数える。ブロックしても期限は延長しない（`./backend/internal/morph/ratelimit/limiter.go` L28-46）。
- **関連**: 本番の秒数は要確認事項 Q-01。

### D-08 必須パラメータとエラー応答（中）
- **ドキュメント**: `./docs/api/openapi.yaml` L27-30 で `html`・`selection`・`auth` を必須とし、L94 で「必須パラメータ欠落」は 400 を返すとしている。
- **コード**:
  - JSON として読めない場合にだけ 400 を返し、必須チェックはしていない（`./backend/function.go` L70-74）。
  - `auth` が欠けると共通鍵が空になるため 401 になる。
  - `selection` や `html` が欠けても 200 を返す。`selection` が空なら `{"tokens": null}` になる（L113-119）。
  - API 仕様にない応答: 405（POST 以外。L63-67）、400 "Invalid HTML"（L97-101）。

### D-09 呼び出し元（Origin）の制限（中）
- **ドキュメント**: `./docs/architecture/03-security.md` L61-67 は、許可する呼び出し元を `chrome-extension://<EXTENSION_ID>` に「厳格に設定」としている。
- **コード**:
  - 環境変数 `ALLOWED_ORIGIN` が未設定または `*` の場合は、チェックを省いてすべて許可し、応答ヘッダーも `*` を返す（`./backend/internal/auth/auth.go` L31-45）。
  - ローカル用の `./docker-compose.yml` は `ALLOWED_ORIGIN` を設定していないため、ローカルでは制限なしで動く。
  - チェックの順序は、本文の読み込みが先で、Origin の検証が後（`./backend/function.go` L70-81）。データフロー図 `./docs/architecture/02-data-flow.mermaid` L12 は Origin の検証を先頭に描いている。

### D-10 README のローカル確認手順（中）
- **ドキュメント**: `./README.md` L19-20 は「バックエンドサーバーを起動した状態」で、`docker-compose.yml` の `SHARED_SECRET=test-secret` を前提としている。起動方法は書いていない。
- **コード・設定**:
  - `make run-backend`（`./Makefile` L17-18）を使うと、`./.env` の値（本番用）を読み込む（L1-2）。共通鍵 `test-secret` や `Origin: chrome-extension://dummy-id` が一致せず、401 や 403 になる見込み。
  - `./Makefile` は `.env` がないとエラーになる。
  - 手順どおりに動くのは `docker compose up` で起動した場合に限られる（`./docker-compose.yml` L14）。
  - 拡張機能側のビルド・読み込み手順（`make dev-extension` など）は README にない。

### D-11 `user_id` の形式（低）
- **ドキュメント**: `./docs/api/openapi.yaml` L58-61 は UUID（v4）としている。`./docs/architecture/03-security.md` L21 は「検証に使用」としている。
- **コード**: 形式を検証せず、空文字も受け付ける（`./backend/function.go` L90）。`./README.md` L33 の例も `test-user-001` で、UUID ではない。
- **補足**: 拡張機能は実際に UUID v4 を生成している（`./extension/src/background.ts` L5-7）。

### D-12 共通鍵の「ハッシュ値」表記（低）
- **ドキュメント**: `./docs/architecture/03-security.md` L10-11 と L33 は「共通のハッシュ値」と照合するとしている。
- **コード**:
  - サーバーはハッシュ化せず、受け取った文字列を環境変数の値とそのまま比較している（`./backend/internal/auth/auth.go` L23-29）。
  - 拡張機能側の値は、ビルド時に環境変数 `VITE_SHARED_SECRET` から JavaScript に埋め込まれる（`./extension/src/utils/config.ts`）。ドキュメントの「設定ファイル（`.env` 等）に保持」とは表現が違う。

### D-13 拡張機能のディレクトリ構成（低）
- **ドキュメント**: `./docs/architecture/04-extension.md` L15-30 は、`utils/` を「共通ロジック（テキスト抽出等）」と説明している。
- **実際**:
  - `utils/` にあるのは `config.ts`（設定値）だけで、テキスト抽出は `content.ts` の中にある。
  - 図にないファイル・フォルダ: `src/vite-env.d.ts`、`.env.*`、`images/`。

### D-14 リクエスト例の項目不足（低）
- **ドキュメント**: `./docs/architecture/02-data-flow.mermaid` L10 と `./docs/architecture/03-security.md` L27-37 のリクエスト例に、`prefix`／`suffix` がない。
- **コード**: 拡張機能はこの 2 項目を送信している（`./extension/src/background.ts` L43-44）。`./docs/api/openapi.yaml` L40-47 にも記載がある。

### D-15 操作の呼び方（低）
- **ドキュメント**: `./docs/architecture/01-overview.md` L9「テキストのドラッグ＆ドロップ（ハイライト）」。
- **コード**: マウスボタンを離した時点（`mouseup`）で選択範囲を取得している（`./extension/src/content.ts` L5）。ドロップ操作は扱っていない。

### D-16 README の「docs/ を正とする」（低・ドキュメント同士）
- **ドキュメント**: `./README.md` L14-15 は「`docs/` 以下のドキュメントを正として開発を進めます」としている。L10 では `docs/` を「設計書、仕様書、アーキテクチャ図」と説明している。
- **食い違う先**: `./CLAUDE.md` は、既存ドキュメントの整合は保証されていないとしている。また、`docs/` は公開ページの置き場所も兼ねている。`notes/` は README の構成説明にない。

### D-17 読みが空文字になる場合（低）
- **ドキュメント**: `./docs/api/openapi.yaml` L87「読みがない場合や記号の場合は空文字またはそのまま」。
- **コード**: 辞書に読みがない場合は表層形（文中の表記そのもの）を返す。空文字は返さない（`./backend/internal/morph/analyzer.go` L36-45）。

### D-18 文脈が一致しないときの代替処理（低）
- **コード**: 選択テキストが前後の文脈と一致する位置に見つからない場合、選択テキストだけを単独で解析して返す（`./backend/function.go` L112-119）。
- **ドキュメント**: `./docs/api/openapi.yaml` など、どのドキュメントにも記載がない。

### C-01〜03 コード内コメントと処理の食い違い（低）
- **C-01** `./backend/function.go` L76: コメントは「Origin 検証（ブラウザからのリクエストの場合）」。実際は、`ALLOWED_ORIGIN` を設定していれば、ブラウザ以外（Origin ヘッダーなし）からのリクエストも 403 で拒否する。
- **C-02** `./backend/function.go` L174-175: 直前の文脈（Prefix）を確認する処理のコメントに、「Suffix チェックする」と書いている。
- **C-03** `./extension/src/content.ts` L62: コメントは `startContainer` だが、実際に使っているのは `commonAncestorContainer`（L63）。

## 要確認事項（コード外の確認が必要）
- **Q-01** 本番の設定値: `RATE_LIMIT_WINDOW_SEC`（D-07）、`ALLOWED_ORIGIN`（D-09）、`VITE_API_ENDPOINT` の URL とパス（D-05）。Git 管理外の秘密ファイルにあるため、今回は開いていない。
- **Q-02** Cloud Functions の世代: `./docs/architecture/01-overview.md` L24 は第 2 世代としている。`./Makefile` のデプロイコマンドには世代を指定する引数がない。実際の世代は GCP（Google Cloud）上で確認が必要。
- **Q-03** コンテナの保存先: `./docs/architecture/01-overview.md` L25 は Artifact Registry としている。コードや設定からは確認できない。
- **Q-04** 精度の主張: `./docs/index.md` L14-15 と L32-34 の読み分けの例（「辛い」「母／父母／伯母」など）が、現行の IPADIC で実際に正しく読めるかは未検証。[読み精度改善の調査](./research/2026-10-11-reading-accuracy-improvement.md)の評価セット作りで、あわせて確認できる。

## 参考（食い違いではないが、今後の作業で注意する点）
- レスポンスの `base_form` という名前は、形態素解析では一般に「原形（活用前の形）」を指す。中身は表層形（文中の表記）で、ドキュメントとコードはこの点で一致している。UniDic へ切り替えるときに取り違えやすい。
- 429（リクエスト過多）用のメッセージは `./extension/src/background.ts` L60-62 に用意してある。ただし画面には出ず、ブラウザの開発者用コンソールに出力されるだけ（`./extension/src/content.ts` L39）。
