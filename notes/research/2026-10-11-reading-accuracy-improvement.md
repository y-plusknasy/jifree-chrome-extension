# 調査: 読み仮名の誤り改善と「意味」表示の実現可能性・コスト

- 調査日: 2026-10-11
- 結論: **軽量ローカル LLM のサーバー配置は非推奨**。まずは辞書を IPADIC → UniDic に切り替えて効果を測る。LLM を使う場合は「辞書＋商用 LLM API」のハイブリッド構成が精度・コストとも有利。
- 追記（2026-10-24）: UniDic への切り替えは実測で読みの向上を確認できなかった（[IPADIC と UniDic の比較](./2026-10-24-ipadic-vs-unidic.md)）。辞書の方針は未決定に戻した。
- 経緯: [セッションメモ](../memo/2026-10-11-session-log.md)

## 1. 現状の仕組みと誤読の原因

- 読みは kagome（Go 製形態素解析器）＋ **IPADIC** の辞書エントリの読み（features[7]）をそのまま返す（[analyzer.go](../../backend/internal/morph/analyzer.go)）。
- 文脈の効き方は「形態素の分割の仕方」のみ。辞書未登録語・分割誤り・多音語（例: 何/人気/生物）で誤読が出る。
- リクエストには `html`（ページ本文）・`selection`・`prefix`/`suffix` が既に含まれる（[types.go](../../backend/internal/api/types.go)）＝文脈を LLM に渡すための材料は現行 API で揃っている。

## 2. 案A: 軽量ローカル LLM をサーバー配置 → 非推奨

### 精度
YOMI-Bench（ANLP 2026 / arXiv 2607.00664）の「読み予測（複数読みを持つ語）」スコア:

| モデル | 規模 | 読み予測（単一読み） | 読み予測（複数読み） |
|---|---|---|---|
| Llama-3.1-Swallow-8B | 8B | 0.82 | 0.51 |
| Llama-3-ELYZA-JP-8B | 8B | 0.79 | 0.44 |
| llm-jp-3-7.2b | 7.2B | 0.80 | 0.52 |
| llm-jp-3-13b | 13B | 0.90 | 0.56 |
| Gemini（2.5 Flash） | 非公開 | 0.98 | 0.86 |
| Claude（Sonnet 4.5） | 非公開 | 1.00 | 0.92 |
| GPT-5 | 非公開 | 0.98 | 0.97 |

- 論文の結論: 日本語特化オープンモデルでも読み性能は低い。サーバーに載る規模では現行辞書より悪化しうる。

### コスト・レイテンシ
- Cloud Run GPU（NVIDIA L4、最小 4 vCPU / 16 GiB）。GPU 単価は第三者集計で約 $0.44〜0.80/時（公式単価は料金ページが動的描画のため未取得）。
- 常時起動: GPU だけで月 $320〜580 程度＋CPU/メモリ。
- ゼロスケール: Cloud Run GPU はインスタンス存続中ずっと課金。コールドスタート（起動約 5 秒＋モデル読み込み数十秒）が操作のたびに発生しうる。

## 3. 案B: 辞書＋商用 LLM API のハイブリッド → LLM 導入時の推奨

### 構成案
1. kagome で分割・辞書読みを算出（現行どおり）
2. 結果と前後文を LLM に渡し、「文脈上の正しい読み」と「子供向けの短い意味」を JSON で返させる
3. LLM 失敗時は辞書結果をそのまま返す（フォールバック）
- UX 案: 辞書結果を即表示 → LLM 結果で上書き（応答は 1〜3 秒程度に伸びる見込み）

### コスト試算（1 回 = 入力約 500 / 出力約 300 トークン と仮定）

| モデル | 単価（入力/出力, $/100万トークン） | 1 回 | 月 1,500 回 |
|---|---|---|---|
| Claude Haiku 5.5 | 0.10 / 0.50 | 約 $0.0002 | 約 $0.3 |
| Gemini 2.5 Flash-Lite | 0.10 / 0.40 | 約 $0.00017 | 約 $0.3 |
| Claude Sonnet 5.5 | 2 / 10 | 約 $0.004 | 約 $6 |

- 思考トークンが加わっても月数ドル規模。現行のレート制限（同一ユーザー 5 秒に 1 回）で暴走も抑止される。
- 単価は 2026-10 時点。改定がありうるため導入前に公式で再確認。

### 注意点
- **プライバシーポリシー改定が必要**: 現行は「第三者提供なし」（[privacy-policy.md](../../docs/privacy-policy.md)）。Gemini 無料枠は送信内容がプロダクト改善に利用されるため、子供の閲覧内容を送るなら有料枠が前提。
- 安価モデル（Haiku / Flash-Lite）自体の読み精度は未測定（上表は上位モデルの値）。

## 4. 追加費用ゼロの改善: IPADIC → UniDic（採用方針 → 2026-10-24 の実測を受けて未決定）

> 追記（2026-10-24）: kagome 標準の UniDic（2.1.2）には仮名形の列がなく、読みの組み立て処理が必要だった。例文 44 問では IPADIC と同点で、メモリは約 4 倍。詳細は [IPADIC と UniDic の比較](./2026-10-24-ipadic-vs-unidic.md)。

- UniDic: 国立国語研究所が現在も人手で保守。語彙数は約 75 万（IPADIC 約 40 万）。kagome は `kagome-dict/uni` で対応。
- 留意点: 分割単位が IPADIC より細かい（短単位）。`findTargetTokens` の選択文字列照合や、フロントのルビ付与単位への影響を確認する必要がある。読みのフィールド位置も IPADIC と異なる。
- 読み精度を IPADIC と直接比較した資料は見つからず → **自前の評価セットで実測**が必要。

## 5. 次の一手

1. 実際に遭遇した誤読例を 20〜50 文集めて評価セット化
2. IPADIC / UniDic（必要なら＋LLM API）で正解率を比較
3. 結果を踏まえて切替を実施し、方針を ADR 化

> 追記（2026-10-24）: 1・2 は自作中心の例文 44 問で IPADIC / UniDic について実施（LLM API は未実施）。切替は見送り、方針は未決定。実際の誤読例の収集は未了。

## 出典
- [YOMI-Bench (arXiv 2607.00664)](https://arxiv.org/abs/2607.00664) / [HTML 版](https://arxiv.org/html/2607.00664) / [ANLP 2026 Q4-12](https://www.anlp.jp/proceedings/annual_meeting/2026/pdf_dir/Q4-12.pdf)
- [Dictionary-Constrained G2P (arXiv 2609.19805)](https://arxiv.org/abs/2609.19805)（辞書ラティス＋LLM生成データで学習、常用漢字読み 99.62%）
- [Gemini API Pricing（公式）](https://ai.google.dev/gemini-api/docs/pricing)
- [Cloud Run GPU（公式）](https://docs.cloud.google.com/run/docs/configuring/services/gpu) / [Cloud Run pricing](https://cloud.google.com/run/pricing)
- [NVIDIA L4 GPU Price Guide 2026](https://jarvislabs.ai/blog/l4-gpu-price) / [2026 GPU Buyer's Guide](https://cerebrium.ai/blog/2026-gpu-buyers-guide)
- [kagome: About the dictionary](https://github.com/ikawaha/kagome/wiki/About-the-dictionary) / [Japanese Tokenizer Dictionaries](https://dampfkraft.com/nlp/japanese-tokenizer-dictionaries.html)
- Claude API 単価: Claude Code 同梱の claude-api リファレンス（2026-10-06 キャッシュ）
