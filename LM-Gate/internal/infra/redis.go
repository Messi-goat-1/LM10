package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisService struct {
	client *redis.Client
}

func NewRedisService(addr string) *RedisService {
	return &RedisService{
		client: redis.NewClient(&redis.Options{
			Addr: addr, // مثلاً "redis:6379"
		}),
	}
}

// 1. تخزين قطع الملف بشكل مؤقت
func (r *RedisService) StoreChunk(fileID string, data []byte) error {
	key := fmt.Sprintf("file:%s:chunks", fileID)
	err := r.client.RPush(ctx, key, data).Err()
	r.client.Expire(ctx, key, 2*time.Hour) // حذف آلي بعد ساعتين لحماية الذاكرة
	return err
}

// 2. تجميع الملف كاملاً عند وصول EOF
func (r *RedisService) GetFullFile(fileID string) ([]byte, error) {
	key := fmt.Sprintf("file:%s:chunks", fileID)
	chunks, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var fullFile []byte
	for _, c := range chunks {
		fullFile = append(fullFile, []byte(c)...)
	}
	return fullFile, nil
}

// 3. فرز البيانات حسب الـ IP (مؤقتاً في Redis قبل MongoDB)
func (r *RedisService) GroupByIP(ip string, data string) error {
	key := fmt.Sprintf("ip_group:%s", ip)
	return r.client.RPush(ctx, key, data).Err()
}

// 4. تنظيف Redis بعد النقل لـ MongoDB
func (r *RedisService) ClearFile(fileID string) {
	r.client.Del(ctx, fmt.Sprintf("file:%s:chunks", fileID))
}

func (r *RedisService) Ping() error {
	return r.client.Ping(ctx).Err()
}
