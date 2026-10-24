# 調査: 辞書 IPADIC と UniDic の比較（kagome でのふりがな用途）

- 調査日: 2026-10-24
- 結論: **kagome 標準の UniDic（kagome-dict/uni）へ切り替えても読みの正しさは向上しなかった**（例文 44 問で IPADIC 35 問・UniDic 補正あり 35 問の同点）。一方でメモリ約 4 倍・起動約 3 倍・読みの組み立て処理が必要になるため、単純な切り替えは利点が乏しい。次回アップデートの方向（下記 6 の A〜C）は未決定。
- 前提の調査: [読み精度改善の調査（2026-10-11）](./2026-10-11-reading-accuracy-improvement.md)

## 1. 比較対象

| | IPADIC | UniDic |
|---|---|---|
| Go モジュール | `github.com/ikawaha/kagome-dict/ipa` v1.2.5（現行） | `github.com/ikawaha/kagome-dict/uni` v1.2.6（最新, 2026-01-03） |
| 収録データ | mecab-ipadic-2.7.0-20070801 | unidic-mecab-2.1.2_src（2013 年）＋ kagome 独自パッチ |
| 解析器 | kagome v2.10.2（共通） | 同左 |

## 2. 辞書の差し替えとライセンス

- **差し替え**: `tokenizer.New(ipa.Dict())` を `tokenizer.New(uni.Dict())` に変えるだけで可能。作者の解説書は辞書のカスタマイズ（再ビルド）手順も示し「UniDic も同様」と明記。ユーザー辞書（`見出し,分割,ヨミ,品詞` 形式、`tokenizer.UserDict`）も辞書種別によらず使える。
- **作者の推奨**: IPADIC 寄り（「メモリが潤沢でない環境では IPADIC」「UniDic は精度が低いという意見が多い」「UniDic は検索の区切りに向く」）。
- **ライセンス**: kagome 本体・kagome-dict/ipa・kagome-dict/uni はいずれも MIT。同梱データは UniDic が BSD 3 条項（UniDic 自体は GPLv2 / LGPLv2.1 / BSD のトリプルライセンスで商用可）、IPADIC が NAIST の条件（再配布時に無保証条項を添付）。サーバー側で動かすだけなら表示義務は原則生じないと読める（法的判断ではない）。

## 3. 読みの項目（Features）の違い

`行った` の `行っ` を解析した結果:

| 辞書 | index | 内容 | 値 |
|---|---|---|---|
| IPADIC | 7 | 読み（現行の [analyzer.go](../../backend/internal/morph/analyzer.go) が使用） | イッ |
| IPADIC | 8 | 発音 | イッ |
| UniDic | 6 | 語彙素読み（lForm） | イク |
| UniDic | 9 | 発音形出現形（pron） | イッ |
| UniDic | 11 | 発音形基本形（pronBase） | イク |

- **UniDic 2.1.2 には「仮名形（kana, 書かれた形の典型的な読み）」の列がない**。元データ自体が 21 列でこの列を含まない。仮名形は UniDic 3.x 系にのみある。
- 語彙素読みは活用前の読み（行っ→イク）、発音形は長音を「ー」で表す（東京→トーキョー）ため、どちらもそのままではふりがなに使えない。
- kagome の `Token.Reading()` は uni では値を返さない（ReadingIndex 未定義）。
- 検証用に「UniDic 補正」を用意: 発音形を基本に、「ー」の位置だけ語彙素読みの同じ位置の文字で置き換える（東京: トーキョー＋トウキョウ→トウキョウ、行っ: イッ のまま）。

## 4. 例文での実測（44 問）

### 方法
- Docker の `golang:1.24` イメージ上で両辞書を動かし、例文中の対象語の読みを比較（プロジェクトのコードは変更せず、一時的な検証プログラムで実施）。
- 対象語の両端が単語の区切りと一致しない場合は、現行の `findTargetTokens` のフォールバックと同じく対象語だけを解析し直す想定で採点（今回は全問で区切りが一致し、フォールバックは 0 件）。
- 例文は多音語・送り仮名・熟字訓などを含む 37 問を作成し、ユーザー提供の例文 1 文から 7 語を加えた。

### 結果

| 読みの取り方 | 正解 |
|---|---|
| IPADIC（index 7） | **35** |
| UniDic 語彙素読み（index 6） | 30 |
| UniDic 発音形（index 9） | 25 |
| UniDic 補正 | **35** |

IPADIC と UniDic 補正で結果が分かれた語:

| 対象（例文） | 正解 | IPADIC | UniDic 補正 |
|---|---|---|---|
| 一日（四月一日は晴れだった。） | ツイタチ | ×イチニチ | ○ |
| 一人（一人で留守番をした。） | ヒトリ | ×イチニン | ○ |
| 父母（ユーザー提供の例文） | フボ | ×チチハハ | ○ |
| 何（何を食べたいですか。） | ナニ | ○ | ×ナン |
| お父さん（お父さんと公園へ行った。） | オトウサン | ○ | ×オチチサン |
| お母さん（ユーザー提供の例文） | オカアサン | ○ | ×オハハサン |

両方とも誤った語: 人気（ヒトケ）、生物（ナマモノ）、色紙（シキシ）、市場（イチバ）、一日中（イチニチジュウ）、二十歳（ハタチ）。文脈で読みが変わる語はどちらの辞書でも区別できない。

両方とも正解した語: 何人、人気（ニンキ）、生物（セイブツ）、上手、下手、東京、日本、大きい、今日、大人、映画、行っ、読ん、通っ（カヨッ / トオッ）、大通り、閉め、牧場、明日、金魚、主人公、大事、竹刀、風邪、上着、果物、私、母、一緒、伯母、後ろ姿、恥ずかし。

- 注意: 例文は少数で大半を自作したため偏りがありうる。実際に遭遇した誤読例を加えると結果が変わる可能性がある。
- 分割の違い: UniDic は「お|母|さん」「主人|公」「二十|歳」のように細かく分ける。フロントは単語ごとにルビを付けるため、表示の単位も細かくなる。

### 負荷

| | IPADIC | UniDic |
|---|---|---|
| 辞書読み込みによるメモリ増分 | 約 93 MB | 約 384 MB |
| 初期化時間 | 約 0.5 秒 | 約 1.5 秒 |
| 単体バイナリサイズ | 15 MB | 47 MB |
| 埋め込み辞書ファイル | 約 11.7 MB | 約 45.5 MB |

- 現行の [Makefile](../../Makefile) は Cloud Functions のメモリを指定していない。既定値（256 MB と認識しているが未確認）のままでは UniDic は載らない見込みで、メモリ設定の引き上げが必要。

## 5. 外部資料

- 小木曽ほか（2010, 言語処理学会）は UniDic・IPAdic・NAIST-jdic のジャンル別比較で UniDic の読み誤りが最も少ないと報告（検索結果の抜粋による数字で本文は未確認。開発元による評価で、辞書は UniDic 1.3 系の時代）。
- それ以外に、第三者による読み正解率の定量比較は見つからなかった。
- UniDic 3.1.0（2021）は仮名形を持つが、辞書データは約 770 MB と大きく、kagome 用の公式ビルドはない（自前ビルドは技術的には可能。Cloud Functions に載るかは未検証）。
- IPADIC は 2007 年頃で更新停止。UniDic は国立国語研究所が保守。

## 6. 次の方向の候補（未決定）

- **A. IPADIC のまま、ユーザー辞書で誤読語を補正**: 低コスト。ただし登録語は文脈によらず常にその読みになるため、一日（ツイタチ / イチニチ）のような語には使えない。
- **B. 文脈依存の誤読に商用 LLM API を併用**: 人気・生物など両辞書とも誤る語に有効な唯一の手段。プライバシーポリシー改定が必要（[前回調査](./2026-10-11-reading-accuracy-improvement.md) の案B）。
- **C. UniDic を採用**: kagome 標準版＋読み補正、または UniDic 3.x の自前ビルド。

## 出典
- [kagome v2 解説書（作者）: 辞書](https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer/viewer/dictionary) / [ユーザー辞書](https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer/viewer/user_dictionary) / [辞書のカスタマイズ](https://zenn.dev/ikawaha/books/kagome-v2-japanese-tokenizer/viewer/custom_dictionary)
- [kagome wiki: About the dictionary](https://github.com/ikawaha/kagome/wiki/About-the-dictionary) / [kagome README](https://github.com/ikawaha/kagome)
- [kagome-dict/uni README](https://github.com/ikawaha/kagome-dict/blob/master/uni/README.md) / [uni/dict.go](https://github.com/ikawaha/kagome-dict/blob/master/uni/dict.go) / [uni/tool/main.go](https://github.com/ikawaha/kagome-dict/blob/master/uni/tool/main.go)
- ライセンス: [kagome LICENSE](https://raw.githubusercontent.com/ikawaha/kagome/v2/LICENSE) / [uni LICENSE](https://raw.githubusercontent.com/ikawaha/kagome-dict/master/uni/LICENSE) / [uni NOTICE](https://raw.githubusercontent.com/ikawaha/kagome-dict/master/uni/NOTICE.txt) / [ipa NOTICE](https://raw.githubusercontent.com/ikawaha/kagome-dict/master/ipa/NOTICE.txt) / [UniDic 商用利用（NINJAL）](https://clrd.ninjal.ac.jp/unidic/en/commerce_use_en.html)
- [unidic-py（読みの各項目の定義）](https://github.com/polm/unidic-py) / [fugashi](https://github.com/polm/fugashi)
- [小木曽ほか 2010「形態素解析辞書のベンチマークテスト」](https://anlp.jp/proceedings/annual_meeting/2010/pdf_dir/PA1-14.pdf)
- [UniDic 旧版一覧（NINJAL）](https://clrd.ninjal.ac.jp/unidic/en/back_number_en.html) / [Japanese Tokenizer Dictionaries](https://www.dampfkraft.com/nlp/japanese-tokenizer-dictionaries.html)
