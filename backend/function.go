package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/y-plusknasy/jifree-chrome-extension/backend/internal/api"
	"github.com/y-plusknasy/jifree-chrome-extension/backend/internal/auth"
	"github.com/y-plusknasy/jifree-chrome-extension/backend/internal/morph"
	"github.com/y-plusknasy/jifree-chrome-extension/backend/internal/ratelimit"
)

var (
	authenticator *auth.Authenticator
	limiter       *ratelimit.Limiter
	analyzer      *morph.Analyzer
)

func init() {
	// 各コンポーネントの初期化
	// Cloud Functionsではinit()で重い処理を行うことで、コールドスタート後の処理を高速化できる
	authenticator = auth.New()
	limiter = ratelimit.New(10 * time.Second) // 10秒制限

	var err error
	analyzer, err = morph.New()
	if err != nil {
		log.Fatalf("failed to initialize analyzer: %v", err)
	}

	// 関数の登録
	functions.HTTP("AnalyzeText", AnalyzeText)
}

// AnalyzeText はHTTPリクエストを処理するエントリーポイントです
func AnalyzeText(w http.ResponseWriter, r *http.Request) {
	// 1. CORSヘッダーの付与
	headers := authenticator.CORSHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	// Preflightリクエスト(OPTIONS)の処理
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// POSTメソッド以外は許可しない
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. リクエストボディの読み込み
	var req api.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Origin検証 (ブラウザからのリクエストの場合)
	origin := r.Header.Get("Origin")
	if !authenticator.ValidateOrigin(origin) {
		http.Error(w, "Forbidden: Invalid Origin", http.StatusForbidden)
		return
	}

	// 4. 共通鍵認証
	if !authenticator.ValidateSecret(req.Auth.SharedSecret) {
		http.Error(w, "Unauthorized: Invalid Secret", http.StatusUnauthorized)
		return
	}

	// 5. レート制限
	if !limiter.Allow(req.Auth.UserID) {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// 6. 形態素分析の実行
	// HTMLからテキストを抽出
	fullText, err := analyzer.ExtractText(req.HTML)
	if err != nil {
		log.Printf("Failed to extract text from HTML: %v", err)
		http.Error(w, "Invalid HTML", http.StatusBadRequest)
		return
	}

	// 全体を解析
	allTokens, err := analyzer.Analyze(fullText)
	if err != nil {
		log.Printf("Analysis failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 7. Selectionの位置特定と抽出
	targetTokens := findTargetTokens(allTokens, req.Selection, req.Prefix, req.Suffix)
	if len(targetTokens) == 0 {
		// 見つからない場合は、フォールバックとしてSelection単体を解析して返す
		// (文脈は考慮されないが、何もしないよりはマシ)
		log.Printf("Target selection not found in context, falling back to simple analysis")
		targetTokens, _ = analyzer.Analyze(req.Selection)
	}

	// 8. レスポンス返却
	resp := api.AnalyzeResponse{
		Tokens: targetTokens,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// findTargetTokens は全トークン列から、指定されたSelectionと文脈(Prefix/Suffix)に一致する部分を探します
func findTargetTokens(tokens []api.Token, selection, prefix, suffix string) []api.Token {
	// Selection自体が空なら何もしない
	if selection == "" {
		return nil
	}

	// 候補を探す
	// Selectionは複数のトークンにまたがる可能性があるため、
	// トークンを連結しながらSelectionと一致するか確認する
	for i := 0; i < len(tokens); i++ {
		// 現在位置からSelectionを構成できるか試行
		var currentSelection string
		var candidateTokens []api.Token

		for j := i; j < len(tokens); j++ {
			currentSelection += tokens[j].BaseForm
			candidateTokens = append(candidateTokens, tokens[j])

			// Selectionと一致したら、文脈チェックへ
			if currentSelection == selection {
				if checkContext(tokens, i, j, prefix, suffix) {
					return candidateTokens
				}
				break // 一致したが文脈が違うなら、この開始位置iは諦めて次へ
			}

			// Selectionより長くなってしまったら、この開始位置iはハズレ
			if len(currentSelection) > len(selection) {
				break
			}
		}
	}

	return nil
}

// checkContext は見つかった候補(startIndex〜endIndex)の前後の文脈が一致するか確認します
func checkContext(tokens []api.Token, startIndex, endIndex int, prefix, suffix string) bool {
	// Prefixチェック
	if prefix != "" {
		// 直前のトークン列を連結して文字列化
		// パフォーマンスのため、直近のトークンから遡って必要な長さだけチェックするのが理想だが、
		// ここでは簡易的に直前すべてを連結してSuffixチェックする
		var prevText string
		// 遡る範囲を制限しても良いが、一旦全部繋げる
		for k := 0; k < startIndex; k++ {
			prevText += tokens[k].BaseForm
		}
		if !strings.HasSuffix(prevText, prefix) {
			return false
		}
	}

	// Suffixチェック
	if suffix != "" {
		var nextText string
		for k := endIndex + 1; k < len(tokens); k++ {
			nextText += tokens[k].BaseForm
		}
		if !strings.HasPrefix(nextText, suffix) {
			return false
		}
	}

	return true
}
