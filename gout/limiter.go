package gout

import (
	"sync"
	"time"
)

type (
	TokenBucket struct {
		capacity  int64      // 桶的容量
		rate      float64    // 令牌放入速率
		tokens    float64    // 当前令牌数量
		lastToken time.Time  // 上一次放令牌的时间
		mtx       sync.Mutex // 互斥锁
	}
)

var (
	TB = &TokenBucket{
		capacity:  0,
		rate:      1.0,
		tokens:    0,
		lastToken: time.Now(),
	}
)

func Limiter(maxConn int64) HandlerFunc {
	TB.capacity = maxConn
	return func(c *Context) {

		if !TB.Allow() {
			c.String(503, "Too many request")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mtx.Lock()
	defer tb.mtx.Unlock()

	now := time.Now()
	// 计算需要放的令牌数量
	tb.tokens = tb.tokens + tb.rate*(float64(now.Sub(tb.lastToken).Seconds())) //上一次的令牌数加上两次请求间的毫秒数
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	// 判断是否允许请求
	if tb.tokens >= 1 {
		tb.tokens--
		tb.lastToken = now
		return true
	} else {
		return false
	}

}
