package ratelimit

import (
	"sync"
	"time"
)

// Limiter はインメモリでレート制限を管理する構造体です。
type Limiter struct {
	store sync.Map
	ttl   time.Duration
}

type entry struct {
	expiresAt time.Time
}

// New は指定されたTTL（有効期間）を持つ新しいLimiterを作成します。
func New(ttl time.Duration) *Limiter {
	return &Limiter{
		ttl: ttl,
	}
}

// Allow は指定されたユーザーIDのリクエストを許可するかどうかを判定します。
// 許可する場合（制限にかかっていない場合）は true を返し、内部状態を更新します。
// 拒否する場合（制限中の場合）は false を返します。
func (l *Limiter) Allow(userID string) bool {
	now := time.Now()

	// 既存のエントリを確認
	if val, ok := l.store.Load(userID); ok {
		e := val.(entry)
		if now.Before(e.expiresAt) {
			// まだ有効期限内なのでブロック
			return false
		}
		// 有効期限切れなので、新しいリクエストとして処理（下で上書き）
	}

	// 新しいアクセスを記録
	l.store.Store(userID, entry{
		expiresAt: now.Add(l.ttl),
	})

	return true
}
