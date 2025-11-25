# アーキテクチャ概要

## システムの目的
まだ漢字が読めない子供向けに、Webページ上の漢字の読み方を提供する。

## コンポーネント

### 1. Chrome Extension (Client)
- **トリガー**: ユーザーによるテキストのドラッグ＆ドロップ（ハイライト）。
- **処理**:
    1. 選択範囲を含むHTML要素を取得。
    2. バックエンドAPIへ送信。
    3. レスポンスを受け取り、UIに読み仮名を表示（ポップアップ等）。

### 2. Google Cloud Functions (Server)
- **ランタイム**: Go
- **処理**:
    1. HTML文字列とハイライト情報を受信。
    2. 形態素分析エンジン（Kagome）を使用して解析。
    3. 漢字に対する読み仮名を特定。
    4. JSON形式でレスポンス。

## インフラストラクチャ
- **Compute**: Google Cloud Functions (Gen 2)
- **Container Registry**: Google Artifact Registry (旧 GCR)
