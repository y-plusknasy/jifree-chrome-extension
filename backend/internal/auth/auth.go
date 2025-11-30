package auth

import (
	"crypto/subtle"
	"os"
)

type Authenticator struct {
	sharedSecret  string
	allowedOrigin string
}

// New は環境変数から設定を読み込み、Authenticatorを初期化します。
func New() *Authenticator {
	return &Authenticator{
		sharedSecret:  os.Getenv("SHARED_SECRET"),
		allowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
	}
}

// ValidateSecret はリクエストされたシークレットが正しいか検証します。
// タイミング攻撃を防ぐため、定数時間比較を使用します。
func (a *Authenticator) ValidateSecret(secret string) bool {
	// 環境変数が設定されていない場合は常に失敗させる（安全側に倒す）
	if a.sharedSecret == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(secret), []byte(a.sharedSecret)) == 1
}

// ValidateOrigin はOriginヘッダーが許可されたものか検証します。
// 開発環境などでALLOWED_ORIGINが未設定、または"*"の場合は、チェックをスキップ（許可）します。
func (a *Authenticator) ValidateOrigin(origin string) bool {
	if a.allowedOrigin == "" || a.allowedOrigin == "*" {
		return true
	}
	return origin == a.allowedOrigin
}

// CORSHeaders はレスポンスに必要なCORSヘッダーを返します。
func (a *Authenticator) CORSHeaders() map[string]string {
	origin := a.allowedOrigin
	if origin == "" {
		origin = "*" // 開発用
	}

	return map[string]string{
		"Access-Control-Allow-Origin":  origin,
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}
}
