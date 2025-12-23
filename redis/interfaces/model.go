package interfaces

import (
	"time"
)

type (
	RedisDbs     map[DatabaseId]DatabaseOps
	DatabaseId   uint
	Keys         map[string]Item
	ExpiringKeys map[string]time.Time
)

// هذه الواجهة هي الأهم، سيستخدمها سيرفر gRPC للوصول للبيانات
type DatabaseOps interface {
	Set(key string, i Item, expires bool, expiry time.Time)
	Get(key string) Item
	Delete(keys ...string) int
	Exists(key string) bool
	Expiry(key string) time.Time
}

// لإدارة حذف البيانات التلقائي
type ExpirationManager interface {
	Start(tick time.Duration, keyNum int, againPercentage int)
	Stop()
}

// الواجهة المبسطة للمحرك الذي سيخدم gRPC
type RedisEngine interface {
	RedisDbs() RedisDbs
	KeyExpirer() ExpirationManager
}
